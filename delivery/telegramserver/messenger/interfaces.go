package messenger

import (
	"context"

	"github.com/go-telegram/bot"
)

type Messenger interface {
	Send(ctx context.Context, params *bot.SendMessageParams) error
	Edit(ctx context.Context, params *bot.EditMessageTextParams) error
	Delete(ctx context.Context, params *bot.DeleteMessageParams) error
}
