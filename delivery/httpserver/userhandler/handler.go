package userhandler

import (
	"telegram-service-platform/service/authservice"
	"telegram-service-platform/service/userservice"
	"telegram-service-platform/service/walletservice"
	"telegram-service-platform/validator/uservalidator"
)

type Handler struct {
	userSvc    *userservice.Service
	walletSvc  *walletservice.Service
	authSvc    *authservice.Service
	authConfig authservice.Config
	botToken   string
	userVal    uservalidator.Validator
}

func New(userSvc *userservice.Service, authSvc *authservice.Service, walletSvc *walletservice.Service,
	userVal uservalidator.Validator, authConfig authservice.Config, botToken string) *Handler {
	return &Handler{
		userSvc:    userSvc,
		authSvc:    authSvc,
		walletSvc:  walletSvc,
		userVal:    userVal,
		authConfig: authConfig,
		botToken:   botToken,
	}
}
