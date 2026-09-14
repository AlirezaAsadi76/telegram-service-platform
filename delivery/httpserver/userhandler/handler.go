package userhandler

import (
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/userservice"
	"telegram-service-platform/service/walletservice"
	"telegram-service-platform/validator/uservalidator"
)

type Handler struct {
	userSvc    *userservice.Service
	walletSvc  *walletservice.Service
	middleware AuthMiddleware
	orderSvc   *orderservice.Service
	botToken   string
	userVal    uservalidator.Validator
}

func New(userSvc *userservice.Service, middleware AuthMiddleware,
	walletSvc *walletservice.Service, orderSvc *orderservice.Service,
	userVal uservalidator.Validator, botToken string) *Handler {
	return &Handler{
		userSvc:    userSvc,
		walletSvc:  walletSvc,
		orderSvc:   orderSvc,
		userVal:    userVal,
		botToken:   botToken,
		middleware: middleware,
	}
}
