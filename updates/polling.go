package updates

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/petuhovskiy/telegram"
)

const delayAfterFail = time.Second * 3
const bufferLength = 50

func StartPolling(bot *telegram.Bot, request telegram.GetUpdatesRequest) (<-chan telegram.Update, error) {
	ch := make(chan telegram.Update, bufferLength)

	go func() {
		for {
			updates, err := bot.GetUpdates(&request)
			if err != nil {
				log.WithError(err).Error("Failed to get updates")
				time.Sleep(delayAfterFail)

				continue
			}

			for _, update := range *updates {
				if update.UpdateID >= request.Offset {
					request.Offset = update.UpdateID + 1
					ch <- update
				}
			}
		}
	}()

	return ch, nil
}
