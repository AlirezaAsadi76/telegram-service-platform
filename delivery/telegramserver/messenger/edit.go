package messenger

import (
	"context"

	"github.com/go-telegram/bot"
)

func (s *Service) Edit(ctx context.Context, params *bot.EditMessageTextParams) error {

	_, err := s.telegramBot.EditMessageText(ctx, params)
	return err
}
