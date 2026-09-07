package mapper

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/params/userparams"

	"github.com/go-telegram/bot/models"
)

func MapTelegramUserToRegisterRequest(tgUser *models.User) userparams.GetOrRegisterRequest {

	return userparams.GetOrRegisterRequest{

		TelegramID: tgUser.ID,
		Username:   tgUser.Username,
		FirstName:  tgUser.FirstName,
		LastName:   tgUser.LastName,
		Role:       entity.UserRole,
	}
}
