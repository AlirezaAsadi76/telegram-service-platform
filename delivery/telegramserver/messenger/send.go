package messenger

import (
	"context"

	"github.com/go-telegram/bot"
)

func (s *Service) Send(ctx context.Context, b *bot.Bot, params *bot.SendMessageParams) error {
	_, err := b.SendMessage(ctx, params)
	return err
}
