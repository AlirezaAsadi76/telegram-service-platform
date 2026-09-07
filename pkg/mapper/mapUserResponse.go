package mapper

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/params/userparams"
)

func MapUserResponse(user *entity.User, isNewUser bool) userparams.GetOrRegisterResponse {
	return userparams.GetOrRegisterResponse{
		UserInfo: userparams.UserInfo{
			Id:         user.ID,
			TelegramID: user.TelegramID,
			Username:   user.Username,
			Role:       user.Role,
		},
		IsNew: isNewUser,
	}

}
