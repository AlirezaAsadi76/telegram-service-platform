package userservice

import (
	"context"
	"fmt"
	"telegram-service-platform/params/userparams"

	"telegram-service-platform/pkg/richerror"
)

func (s Service) FindUserByTelegramID(ctx context.Context, req userparams.FindUserByTelegramIDRequest) (userparams.FindUserByTelegramIDResponse, error) {
	const op = "userservice.FindUserByTelegramID"

	user, err := s.repository.FindUserByTelegramID(ctx, req.TelegramID.Int64())
	fmt.Println(user, err)
	if err != nil {
		if richerror.IsKind(err, richerror.KindNotFound) {
			return userparams.FindUserByTelegramIDResponse{Found: false}, nil
		}
		return userparams.FindUserByTelegramIDResponse{}, richerror.New(op, err)
	}

	return userparams.FindUserByTelegramIDResponse{
		UserInfo: userparams.UserInfo{
			Username:   user.Username,
			TelegramID: user.TelegramID,
			Id:         user.ID,
			Role:       user.Role,
		},
		Found: true,
	}, nil
}
