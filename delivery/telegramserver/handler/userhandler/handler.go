package userhandler

import (
	"telegram-service-platform/delivery/telegramserver/messenger"
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/userservice"
	"telegram-service-platform/service/walletservice"
	"telegram-service-platform/validator/uservalidator"
)

type Handler struct {
	userService   *userservice.Service
	userValidator uservalidator.Validator
	walletService *walletservice.Service
	orderService  *orderservice.Service
	messenger     messenger.Messenger
}

func New(userService *userservice.Service, walletService *walletservice.Service,
	orderService *orderservice.Service, userValidator uservalidator.Validator, messenger messenger.Messenger) Handler {
	return Handler{
		userService:   userService,
		walletService: walletService,
		orderService:  orderService,
		userValidator: userValidator,
		messenger:     messenger,
	}
}
