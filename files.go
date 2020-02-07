package telegram

import (
	"fmt"
)

func (b *Bot) GetFileURL(fileID string, uploader func(tempURL string) (string, error)) (string, error) {
	file, err := b.GetFile(&GetFileRequest{
		FileID: fileID,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/file/bot%s/%s", requestAddress, b.token, file.FilePath)
	return uploader(url)
}
