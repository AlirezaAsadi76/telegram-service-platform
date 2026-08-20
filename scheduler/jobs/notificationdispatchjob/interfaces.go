package notificationdispatchjob

import (
	"context"
	"time"

	"github.com/go-telegram/bot"
)

// TelegramSender interface minimal برای جلوگیری از وابستگی سخت به bot
type TelegramSender interface {
	Send(ctx context.Context, params *bot.SendMessageParams) error
}

type RedisRepository interface {
	BRPop(ctx context.Context, ttl time.Duration, key string) ([]string, error)
}
