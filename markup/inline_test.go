package markup

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/petuhovskiy/telegram"
)

// TestMarshalRequest tests that empty ReplyMarkup will be not included in the request.
func TestMarshalRequest(t *testing.T) {
	req := telegram.SendMessageRequest{
		ChatID:      "123",
		Text:        "123",
		ReplyMarkup: nil,
	}

	j, err := json.Marshal(req)
	assert.Nil(t, err)
	assert.Equal(t, `{"chat_id":"123","text":"123"}`, string(j))
}

// TestMarshalRequest2 tests that empty ReplyMarkup will be not included in the request.
func TestMarshalRequest2(t *testing.T) {
	req := telegram.SendMessageRequest{
		ChatID:      "123",
		Text:        "123",
		ReplyMarkup: InlineKeyboard(nil),
	}

	j, err := json.Marshal(req)
	assert.Nil(t, err)
	assert.Equal(t, `{"chat_id":"123","text":"123"}`, string(j))
}
