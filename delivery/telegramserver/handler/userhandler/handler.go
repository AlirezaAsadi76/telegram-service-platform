package userhandler

import (
	"telegram-service-platform/delivery/telegramserver/messenger"
	"telegram-service-platform/service/userservice"
	"telegram-service-platform/validator/uservalidator"
)

type Handler struct {
	userService   *userservice.Service
	userValidator uservalidator.Validator
	messenger     messenger.Messenger
}

func New(userService *userservice.Service, userValidator uservalidator.Validator, messenger messenger.Messenger) Handler {
	return Handler{
		userService:   userService,
		userValidator: userValidator,
		messenger:     messenger,
	}
}
