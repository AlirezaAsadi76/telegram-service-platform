package orderhandler

import (
	"telegram-service-platform/delivery/telegramserver/messenger"
	"telegram-service-platform/service/checkoutservice"
	"telegram-service-platform/service/orderflowservice"
	"telegram-service-platform/service/pricingservice"
	"telegram-service-platform/service/productservice"
	"telegram-service-platform/service/userservice"
	"telegram-service-platform/validator/ordervalidator"
)

type Handler struct {
	productService   *productservice.Service
	checkoutService  *checkoutservice.Service
	orderFlowService *orderflowservice.Service
	userService      *userservice.Service
	pricingSvc       *pricingservice.Service
	messenger        messenger.Messenger
	validator        ordervalidator.Validator
}

func New(
	productService *productservice.Service,
	checkoutService *checkoutservice.Service,
	orderFlowService *orderflowservice.Service,
	pricingSvc *pricingservice.Service,
	userService *userservice.Service,
	messenger messenger.Messenger,
	validator ordervalidator.Validator,

) *Handler {
	return &Handler{
		productService:   productService,
		checkoutService:  checkoutService,
		orderFlowService: orderFlowService,
		messenger:        messenger,
		pricingSvc:       pricingSvc,
		userService:      userService,
		validator:        validator,
	}
}
