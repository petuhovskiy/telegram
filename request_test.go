package telegram

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockUploader struct{}

func (m mockUploader) Name() string {
	panic("implement me")
}

func (m mockUploader) Reader() (io.Reader, error) {
	panic("implement me")
}

func (m mockUploader) Size() int64 {
	panic("implement me")
}

func TestIsFileUpload(t *testing.T) {
	check := func(req interface{}, expected fileUpload) {
		actual, ok := isFileUpload(req)

		assert.Equal(t, expected.fieldname != "", ok)
		assert.Equal(t, expected, actual)
	}

	check(&SendMessageRequest{}, fileUpload{})
	check(&SendAnimationRequest{}, fileUpload{})
	check(&SendPhotoRequest{}, fileUpload{})
	check(&SendPhotoRequest{
		ChatID:              "12345",
		Photo:               mockUploader{},
		Caption:             "caption here",
		ParseMode:           "html",
		DisableNotification: true,
		ReplyToMessageID:    101,
		ReplyMarkup: InlineKeyboardMarkup{
			InlineKeyboard: [][]InlineKeyboardButton{{{
				Text: "hello",
			}}},
		},
	}, fileUpload{
		params: map[string]string{
			"chat_id":              "12345",
			"caption":              "caption here",
			"parse_mode":           "html",
			"disable_notification": "true",
			"reply_to_message_id":  "101",
			"reply_markup":         `{"inline_keyboard":[[{"text":"hello"}]]}`,
		},
		fieldname: "photo",
		file:      mockUploader{},
		err:       nil,
	})
	check(&SendPhotoRequest{
		Photo: mockUploader{},
	}, fileUpload{
		params: map[string]string{
			"chat_id": "",
		},
		fieldname: "photo",
		file:      mockUploader{},
		err:       nil,
	})
}
