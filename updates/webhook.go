package updates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/petuhovskiy/telegram"
)

const defaultMaxBytes = 1 << 20

type Opts struct {
	CertFile     string
	KeyFile      string
	AddrWh       string
	Salt         string
	Timeout      int
	Port         string // ":80"
	HandleUpdate func(update *telegram.Update)
}

type Webhook struct {
	opts    *Opts
	bot     *telegram.Bot
	botHash string
	address *url.URL
	server  *http.Server
}

func NewWebhook(bot *telegram.Bot, opts *Opts) (*Webhook, error) {
	wh := &Webhook{
		opts:    opts,
		bot:     bot,
		botHash: telegram.GetHash(opts.Salt, bot),
	}

	if _, err := os.Stat(opts.CertFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("cert file error: %w", err)
	}
	if _, err := os.Stat(opts.KeyFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("key file error: %w", err)
	}

	var err error
	wh.address, err = url.Parse(opts.AddrWh)
	if err != nil {
		return nil, err
	}

	addrWH := wh.address
	addrWH.Path = path.Join(addrWH.Path, wh.botHash)

	_, err = bot.SetWebhook(&telegram.SetWebhookRequest{
		URL: addrWH.String(),
	})
	if err != nil {
		return nil, err
	}

	// Init and run HTTP Server for Webhook
	wh.server = &http.Server{
		Addr:           opts.Port,
		ReadTimeout:    time.Duration(opts.Timeout) * time.Second,
		WriteTimeout:   time.Duration(opts.Timeout) * time.Second,
		MaxHeaderBytes: defaultMaxBytes,
		Handler:        wh,
	}

	return wh, nil
}

func (wh *Webhook) Start() error {
	log.WithField("port", wh.server.Addr).Info("Starting http server")
	if err := wh.server.ListenAndServeTLS(wh.opts.CertFile, wh.opts.KeyFile); err != nil {
		log.WithError(err).Error("failed to ListenAndServeTLS")
		return err
	}
	return nil
}

func (wh *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.WithField("url", r.URL.String()).Debug("ServeHTTP")

	if strings.HasPrefix(r.URL.Path, "/wh/") {
		args := strings.Split(strings.Trim(r.URL.Path, "/ "), "/")
		mHash := args[len(args)-1]

		// Находим по хэшу нашего бота
		if mHash != wh.botHash {
			log.Error("webhook bot hash mismatch")
			return
		}

		upd, err := parseRequest(r)
		if err != nil {
			log.WithField("err", err).Error("failed to parse telegram webhook")
			return
		}

		wh.opts.HandleUpdate(upd)
	}

	w.WriteHeader(http.StatusOK)
}

func parseRequest(r *http.Request) (*telegram.Update, error) {
	defer r.Body.Close()
	var upd telegram.Update

	err := json.NewDecoder(r.Body).Decode(&upd)
	if err != nil {
		return nil, err
	}

	return &upd, nil
}
