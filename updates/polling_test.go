package updates

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/petuhovskiy/telegram"
)

func TestStartPolling(t *testing.T) {
	token := os.Getenv("TELEGRAM_BOT_TEST_TOKEN")
	if token == "" {
		t.SkipNow()
	}

	bot := telegram.NewBot(token)
	ch, err := StartPolling(bot, telegram.GetUpdatesRequest{
		Offset: 0,
	})
	if err != nil {
		t.Fatal(err)
	}

	update := <-ch
	spew.Dump(update)
	// Output: update sent to test bot
}
