package messenger

import (
	"context"

	"github.com/go-telegram/bot"
)

func (s *Service) Delete(ctx context.Context, b *bot.Bot, params *bot.DeleteMessageParams) error {

	_, err := b.DeleteMessage(ctx, params)

	return err
}
