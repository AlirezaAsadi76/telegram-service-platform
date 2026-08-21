package userservice

import (
	"context"
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params"
	"telegram-service-platform/pkg/mapper"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (s Service) GetOrRegister(ctx context.Context, request params.GetOrRegisterRequest) (params.GetOrRegisterResponse, error) {
	const Op = "userservice.GetOrRegister"

	existingUser, fErr := s.repository.FindUserByTelegramID(ctx, request.TelegramID)

	if fErr != nil && !richerror.IsKind(fErr, richerror.KindNotFound) {
		return params.GetOrRegisterResponse{}, richerror.New(Op, fErr)
	}
	if existingUser != nil {
		return mapper.MapUserResponse(existingUser, false), nil
	}
	fmt.Println("user : ", request.TelegramID)
	user := entity.User{
		Username:   request.Username,
		TelegramID: request.TelegramID,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		Role:       request.Role,
	}

	rErr := s.repository.Create(ctx, &user)
	if rErr != nil {
		return params.GetOrRegisterResponse{},
			richerror.New(Op, rErr).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.InternalServerError)
	}

	if _, gcErr := s.walletSvc.GetOrCreate(ctx, user.ID, entity.CurrencyTOMAN); gcErr != nil {

		logger.Logger.Error("ensure wallet failed",
			zap.String("op", Op),
			zap.Uint64("user_id", user.ID),
			zap.Error(gcErr),
		)
		return params.GetOrRegisterResponse{},
			richerror.New(Op, rErr).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.InternalServerError)
	}
	return mapper.MapUserResponse(&user, true), nil
}
