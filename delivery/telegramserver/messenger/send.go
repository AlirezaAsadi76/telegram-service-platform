package messenger

import (
	"context"

	"github.com/go-telegram/bot"
)

func (s *Service) Send(ctx context.Context, params *bot.SendMessageParams) error {
	_, err := s.telegramBot.SendMessage(ctx, params)
	return err
}
