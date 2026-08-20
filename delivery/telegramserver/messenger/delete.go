package messenger

import (
	"context"

	"github.com/go-telegram/bot"
)

func (s *Service) Delete(ctx context.Context, params *bot.DeleteMessageParams) error {

	_, err := s.telegramBot.DeleteMessage(ctx, params)

	return err
}
