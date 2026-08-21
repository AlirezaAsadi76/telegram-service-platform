package mapper

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/params"
)

func MapUserResponse(user *entity.User, isNewUser bool) params.GetOrRegisterResponse {
	return params.GetOrRegisterResponse{
		UserInfo: params.UserInfo{
			Id:         user.ID,
			TelegramID: user.TelegramID,
			Username:   user.Username,
			Role:       user.Role,
		},
		IsNew: isNewUser,
	}

}
