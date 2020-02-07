package markup

import (
	"github.com/petuhovskiy/telegram"
)

func InlineKeyboardRows(maxCountInRow int, buttons []telegram.InlineKeyboardButton) [][]telegram.InlineKeyboardButton {
	var keyboard [][]telegram.InlineKeyboardButton

	for i := 0; i < len(buttons); i += maxCountInRow {
		if i+maxCountInRow >= len(buttons) {
			keyboard = append(
				keyboard,
				buttons[i:],
			)
		} else {
			keyboard = append(
				keyboard,
				buttons[i:i+maxCountInRow],
			)
		}
	}

	return keyboard
}

func AppendRowButton(keyboard [][]telegram.InlineKeyboardButton, button telegram.InlineKeyboardButton) [][]telegram.InlineKeyboardButton {
	return append(
		keyboard,
		[]telegram.InlineKeyboardButton{button},
	)
}
