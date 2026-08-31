package keyboard

import "github.com/go-telegram/bot/models"

type ButtonStyle string

const (
	Primary ButtonStyle = "primary"
	Success ButtonStyle = "success"
	Danger  ButtonStyle = "danger"
)

type Button struct {
	Text  string
	Data  string
	URL   string
	Style ButtonStyle
	Icon  string
}

func (b Button) ToTelegramInlineButton() models.InlineKeyboardButton {
	return models.InlineKeyboardButton{
		Text:              b.Text,
		CallbackData:      b.Data,
		Style:             string(b.Style),
		URL:               b.URL,
		IconCustomEmojiID: b.Icon,
	}
}
