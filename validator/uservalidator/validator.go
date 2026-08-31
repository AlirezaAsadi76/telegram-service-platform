package uservalidator

import "telegram-service-platform/service/userservice"

type Validator struct {
	userService *userservice.Service
}

func New(userService *userservice.Service) Validator {
	return Validator{
		userService: userService,
	}
}
