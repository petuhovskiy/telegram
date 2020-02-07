package markup

import (
	"github.com/petuhovskiy/telegram"
)

func Pagination(pageSize, rowSize, length, page int, setPage func(int) string, elem func(int) telegram.InlineKeyboardButton) [][]telegram.InlineKeyboardButton {
	pages := (length + pageSize - 1) / pageSize

	var buttons []telegram.InlineKeyboardButton
	for i := page * pageSize; i < (page+1)*pageSize && i < length; i++ {
		buttons = append(buttons, elem(i))
	}

	var arrows []telegram.InlineKeyboardButton
	if page > 0 {
		arrows = append(arrows, telegram.InlineKeyboardButton{
			Text:         "<-",
			CallbackData: setPage(page - 1),
		})
	}

	if page+1 < pages {
		arrows = append(arrows, telegram.InlineKeyboardButton{
			Text:         "->",
			CallbackData: setPage(page + 1),
		})
	}

	keyboard := InlineKeyboardRows(rowSize, buttons)

	if len(arrows) != 0 {
		keyboard = append(keyboard, arrows)
	}

	return keyboard
}

func Pagination2(rowSize, page int, hasNextPage bool, setPage func(int) string, buttons []telegram.InlineKeyboardButton) [][]telegram.InlineKeyboardButton {
	var arrows []telegram.InlineKeyboardButton
	if page > 0 {
		arrows = append(arrows, telegram.InlineKeyboardButton{
			Text:         "<-",
			CallbackData: setPage(page - 1),
		})
	}

	if hasNextPage {
		arrows = append(arrows, telegram.InlineKeyboardButton{
			Text:         "->",
			CallbackData: setPage(page + 1),
		})
	}

	keyboard := InlineKeyboardRows(rowSize, buttons)

	if len(arrows) != 0 {
		keyboard = append(keyboard, arrows)
	}

	return keyboard
}
