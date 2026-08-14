package notificationservice

import (
	"context"
	"telegram-service-platform/entity/notificationentity"
)

type Repository interface {
	Create(ctx context.Context, notification *notificationentity.Notification) error
	GetPending(ctx context.Context, limit int) ([]notificationentity.Notification, error)
	UpdateStatus(ctx context.Context, id uint64, status notificationentity.NotificationStatus) error
	UpdateRetryCount(ctx context.Context, id uint64, retryCount int) error
}

type RedisRepository interface {
	LPush(ctx context.Context, queueKey string, data any) error
}
