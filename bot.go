package telegram

import (
	"net/http"
)

type Opts struct {
	Client     *http.Client
	Middleware func(RequestHandler) RequestHandler
}

type Bot struct {
	makeRequest RequestHandler
	token       string
	client      *http.Client
}

func NewBot(token string) *Bot {
	return NewBotWithOpts(token, nil)
}

func NewBotWithOpts(token string, opts *Opts) *Bot {
	b := &Bot{
		token:  token,
		client: http.DefaultClient,
	}
	b.makeRequest = b.executeRequest

	if opts != nil {
		if opts.Client != nil {
			b.client = opts.Client
		}
		if opts.Middleware != nil {
			b.makeRequest = opts.Middleware(b.makeRequest)
		}
	}

	return b
}
