package notificationdispatchjob

import (
	"context"
	"time"
)

// TelegramSender interface minimal برای جلوگیری از وابستگی سخت به bot
type TelegramSender interface {
	SendText(ctx context.Context, chatID int64, text string) error
}

type RedisRepository interface {
	BRPop(ctx context.Context, ttl time.Duration, key string) ([]string, error)
}
