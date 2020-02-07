package apigen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitLongLine(t *testing.T) {
	line := "Use this method to send audio files, if you want Telegram clients to display them in the music player. Your audio must be in the .MP3 or .M4A format. On success, the sent Message is returned. Bots can currently send audio files of up to 50 MB in size, this limit may be changed in the future."

	splitted := splitLongComment(line)

	expected := []string{
		"Use this method to send audio files, if you want Telegram clients to display",
		"them in the music player. Your audio must be in the .MP3 or .M4A format. On",
		"success, the sent Message is returned. Bots can currently send audio files of up",
		"to 50 MB in size, this limit may be changed in the future.",
	}

	assert.Equal(t, expected, splitted)
}
