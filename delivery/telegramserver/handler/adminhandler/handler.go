package adminhandler

import (
	"telegram-service-platform/config"
	"telegram-service-platform/delivery/telegramserver/messenger"
	"telegram-service-platform/service/checkoutservice"
	"telegram-service-platform/service/userservice"
)

type Handler struct {
	checkoutService *checkoutservice.Service
	userService     *userservice.Service
	messenger       messenger.Messenger
	config          config.AdminConfig
}

func New(
	checkoutService *checkoutservice.Service,
	userService *userservice.Service,
	messenger messenger.Messenger,
	config config.AdminConfig,
) *Handler {
	return &Handler{
		checkoutService: checkoutService,
		userService:     userService,
		messenger:       messenger,
		config:          config,
	}
}
