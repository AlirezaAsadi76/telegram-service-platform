package messenger

import (
	"context"

	"github.com/go-telegram/bot"
)

func (s *Service) Edit(ctx context.Context, b *bot.Bot, params *bot.EditMessageTextParams) error {

	_, err := b.EditMessageText(ctx, params)
	return err
}
