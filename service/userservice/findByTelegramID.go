package userservice

import (
	"context"
	"fmt"

	"telegram-service-platform/params"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) FindUserByTelegramID(ctx context.Context, req params.FindUserByTelegramIDRequest) (params.FindUserByTelegramIDResponse, error) {
	const op = "userservice.FindUserByTelegramID"

	user, err := s.repository.FindUserByTelegramID(ctx, req.TelegramID.Int64())
	fmt.Println(user, err)
	if err != nil {
		if richerror.IsKind(err, richerror.KindNotFound) {
			return params.FindUserByTelegramIDResponse{Found: false}, nil
		}
		return params.FindUserByTelegramIDResponse{}, richerror.New(op, err)
	}

	return params.FindUserByTelegramIDResponse{
		UserInfo: params.UserInfo{
			Username:   user.Username,
			TelegramID: user.TelegramID,
			Id:         user.ID,
			Role:       user.Role,
		},
		Found: true,
	}, nil
}
