package userhandler

import (
	"telegram-service-platform/service/userservice"
	"telegram-service-platform/validator/uservalidator"
)

type Handler struct {
	userService   userservice.Service
	userValidator uservalidator.Validator
}

func New(userService userservice.Service, userValidator uservalidator.Validator) Handler {
	return Handler{
		userService:   userService,
		userValidator: userValidator,
	}
}
