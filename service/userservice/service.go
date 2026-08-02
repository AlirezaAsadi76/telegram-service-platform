package userservice

import (
	"context"
	"telegram-service-platform/service/entity"
)

type Repository interface {
	Register(ctx context.Context, user entity.User) (entity.User, error)
}

type Service struct {
	repository Repository
}

func New(repository Repository) Service {
	return Service{
		repository: repository,
	}
}
