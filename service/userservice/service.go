package userservice

import (
	"context"
	"telegram-service-platform/entity"
)

type UserRepository interface {
	FindUserByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}

type Service struct {
	repository UserRepository
}

func New(repository UserRepository) *Service {
	return &Service{
		repository: repository,
	}
}
