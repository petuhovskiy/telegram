package telegram

// This object represents one result of an inline query. Telegram clients currently
// support results of the following 20 types:
// - InlineQueryResultCachedAudio
// - InlineQueryResultCachedDocument
// - InlineQueryResultCachedGif
// - InlineQueryResultCachedMpeg4Gif
// - InlineQueryResultCachedPhoto
// - InlineQueryResultCachedSticker
// - InlineQueryResultCachedVideo
// - InlineQueryResultCachedVoice
// - InlineQueryResultArticle
// - InlineQueryResultAudio
// - InlineQueryResultContact
// - InlineQueryResultGame
// - InlineQueryResultDocument
// - InlineQueryResultGif
// - InlineQueryResultLocation
// - InlineQueryResultMpeg4Gif
// - InlineQueryResultPhoto
// - InlineQueryResultVenue
// - InlineQueryResultVideo
// - InlineQueryResultVoice
type InlineQueryResult interface{}

// This object represents the content of a message to be sent as a result of an
// inline query. Telegram clients currently support the following 4 types:
// - InputTextMessageContent
// - InputLocationMessageContent
// - InputVenueMessageContent
// - InputContactMessageContent
type InputMessageContent interface{}
