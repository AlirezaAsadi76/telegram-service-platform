package messenger

import (
	"context"

	"github.com/go-telegram/bot"
)

type Messenger interface {
	Send(ctx context.Context, b *bot.Bot, params *bot.SendMessageParams) error
	Edit(ctx context.Context, b *bot.Bot, params *bot.EditMessageTextParams) error
	Delete(ctx context.Context, b *bot.Bot, params *bot.DeleteMessageParams) error
}
