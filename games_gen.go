// Code generated by telegram-apigen. DO NOT EDIT.

package telegram

import "encoding/json"

type SendGameRequest struct {
	// Unique identifier for the target chat
	ChatID int `json:"chat_id"`

	// Short name of the game, serves as the unique identifier for the game. Set up
	// your games via Botfather.
	GameShortName string `json:"game_short_name"`

	// Optional. Sends the message silently. Users will receive a notification with no
	// sound.
	DisableNotification bool `json:"disable_notification,omitempty"`

	// Optional. If the message is a reply, ID of the original message
	ReplyToMessageID int `json:"reply_to_message_id,omitempty"`

	// Optional. Pass True, if the message should be sent even if the specified
	// replied-to message is not found
	AllowSendingWithoutReply bool `json:"allow_sending_without_reply,omitempty"`

	// Optional. A JSON-serialized object for an inline keyboard. If empty, one 'Play
	// game_title' button will be shown. If not empty, the first button must launch the
	// game.
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// Use this method to send a game. On success, the sent Message is returned.
func (b *Bot) SendGame(req *SendGameRequest) (*Message, error) {
	j, err := b.makeRequest("sendGame", req)
	if err != nil {
		return nil, err
	}

	var resp Message
	err = json.Unmarshal(j, &resp)
	return &resp, err
}

// This object represents a game. Use BotFather to create and edit games, their
// short names will act as unique identifiers.
type Game struct {
	// Title of the game
	Title string `json:"title"`

	// Description of the game
	Description string `json:"description"`

	// Photo that will be displayed in the game message in chats.
	Photo []PhotoSize `json:"photo"`

	// Optional. Brief description of the game or high scores included in the game
	// message. Can be automatically edited to include current high scores for the game
	// when the bot calls setGameScore, or manually edited using editMessageText.
	// 0-4096 characters.
	Text string `json:"text,omitempty"`

	// Optional. Special entities that appear in text, such as usernames, URLs, bot
	// commands, etc.
	TextEntities []MessageEntity `json:"text_entities,omitempty"`

	// Optional. Animation that will be displayed in the game message in chats. Upload
	// via BotFather
	Animation *Animation `json:"animation,omitempty"`
}

type SetGameScoreRequest struct {
	// User identifier
	UserID int `json:"user_id"`

	// New score, must be non-negative
	Score int `json:"score"`

	// Optional. Pass True, if the high score is allowed to decrease. This can be
	// useful when fixing mistakes or banning cheaters
	Force bool `json:"force,omitempty"`

	// Optional. Pass True, if the game message should not be automatically edited to
	// include the current scoreboard
	DisableEditMessage bool `json:"disable_edit_message,omitempty"`

	// Optional. Required if inline_message_id is not specified. Unique identifier for
	// the target chat
	ChatID int `json:"chat_id,omitempty"`

	// Optional. Required if inline_message_id is not specified. Identifier of the sent
	// message
	MessageID int `json:"message_id,omitempty"`

	// Optional. Required if chat_id and message_id are not specified. Identifier of
	// the inline message
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

// Use this method to set the score of the specified user in a game. On success, if
// the message was sent by the bot, returns the edited Message, otherwise returns
// True. Returns an error, if the new score is not greater than the user's current
// score in the chat and force is False.
func (b *Bot) SetGameScore(req *SetGameScoreRequest) (*Message, error) {
	j, err := b.makeRequest("setGameScore", req)
	if err != nil {
		return nil, err
	}

	var resp Message
	err = json.Unmarshal(j, &resp)
	return &resp, err
}

type GetGameHighScoresRequest struct {
	// Target user id
	UserID int `json:"user_id"`

	// Optional. Required if inline_message_id is not specified. Unique identifier for
	// the target chat
	ChatID int `json:"chat_id,omitempty"`

	// Optional. Required if inline_message_id is not specified. Identifier of the sent
	// message
	MessageID int `json:"message_id,omitempty"`

	// Optional. Required if chat_id and message_id are not specified. Identifier of
	// the inline message
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

// Use this method to get data for high score tables. Will return the score of the
// specified user and several of their neighbors in a game. On success, returns an
// Array of GameHighScore objects.
//
//
// This method will currently return scores for the target user, plus two of their
// closest neighbors on each side. Will also return the top three users if the user
// and his neighbors are not among them. Please note that this behavior is subject
// to change.
//
func (b *Bot) GetGameHighScores(req *GetGameHighScoresRequest) (*GameHighScore, error) {
	j, err := b.makeRequest("getGameHighScores", req)
	if err != nil {
		return nil, err
	}

	var resp GameHighScore
	err = json.Unmarshal(j, &resp)
	return &resp, err
}

// This object represents one row of the high scores table for a game.
//
// And that's about all we've got for now.
// If you've got any questions, please check out our Bot FAQ »
type GameHighScore struct {
	// Position in high score table for the game
	Position int `json:"position"`

	// User
	User *User `json:"user"`

	// Score
	Score int `json:"score"`
}
