package userservice

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/params"
	"telegram-service-platform/pkg/mapper"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetOrRegister(ctx context.Context, request params.GetOrRegisterRequest) (params.GetOrRegisterResponse, error) {
	const Op = "userservice.GetOrRegister"

	existingUser, fErr := s.repository.FindUserByTelegramID(ctx, request.TelegramID)
	if fErr != nil {
		return params.GetOrRegisterResponse{}, richerror.New(Op, fErr)
	}
	if existingUser != nil {
		return mapper.MapUserResponse(existingUser), nil
	}

	user := entity.User{
		Username:   request.Username,
		TelegramID: request.TelegramID,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		Role:       request.Role,
	}

	userRegistered, rErr := s.repository.Create(ctx, &user)
	if rErr != nil {
		return params.GetOrRegisterResponse{},
			richerror.New(Op, rErr).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.InternalServerError)
	}

	return mapper.MapUserResponse(userRegistered), nil
}
