package userservice

import (
	"context"
	"telegram-service-platform/service/entity"
	"telegram-service-platform/service/params"
	"telegram-service-platform/service/pkg/richerror"
)

func (s Service) Register(ctx context.Context, request params.RegisterRequest) (params.RegisterResponse, error) {
	const Op = "userservice.Register"

	user := entity.User{
		Username: request.Username,
		TelegramID: request.TelegramID,
		FirstName: request.FirstName,
		LastName: request.LastName,
		Role: request.Role,
	}

	userRegisted, rErr := s.repository.Register(ctx, user)
	if rErr != nil {
		return params.RegisterResponse{}, richerror.New(Op,rErr).
	}
	return params.RegisterResponse{}, nil
}
