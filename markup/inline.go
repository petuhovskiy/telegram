package markup

import "github.com/petuhovskiy/telegram"

func InlineKeyboard(keyboard [][]telegram.InlineKeyboardButton) telegram.AnyKeyboard {
	if len(keyboard) == 0 {
		return nil
	}

	return &telegram.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

func InlineKeyboardMarkup(keyboard [][]telegram.InlineKeyboardButton) *telegram.InlineKeyboardMarkup {
	if len(keyboard) == 0 {
		return nil
	}

	return &telegram.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}
