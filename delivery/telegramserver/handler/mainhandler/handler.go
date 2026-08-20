package mainhandler

import (
	"telegram-service-platform/delivery/telegramserver/messenger"
	"telegram-service-platform/service/productservice"
	"telegram-service-platform/service/userservice"
)

type Handler struct {
	productService *productservice.Service
	userService    *userservice.Service
	messenger      messenger.Messenger
}

func New(
	productService *productservice.Service,
	userService *userservice.Service,
	messenger messenger.Messenger,
) *Handler {
	return &Handler{
		productService: productService,
		userService:    userService,
		messenger:      messenger,
	}
}
