package mapper

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/params"

	"github.com/go-telegram/bot/models"
)

func MapTelegramUserToRegisterRequest(tgUser *models.User) params.GetOrRegisterRequest {

	return params.GetOrRegisterRequest{

		TelegramID: tgUser.ID,
		Username:   tgUser.Username,
		FirstName:  tgUser.FirstName,
		LastName:   tgUser.LastName,
		Role:       entity.UserRole,
	}
}
