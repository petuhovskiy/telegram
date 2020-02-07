package telegram

import (
	"crypto/sha256"
	"fmt"
)

func GetHash(salt string, bot *Bot) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(salt+bot.token+salt)))
}
