package userservice

import (
	"context"
	"telegram-service-platform/entity"
)

type UserRepository interface {
	FindUserByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	UpdateLastSeenBulk(ctx context.Context, telegramIDs []int64) (int64, error)
}

type ActivityTrackerRepository interface {
	TrackActivity(ctx context.Context, telegramID int64) error
	GetActiveUsers(ctx context.Context) ([]int64, error)
	ClearActivity(ctx context.Context, telegramID int64) error
}
