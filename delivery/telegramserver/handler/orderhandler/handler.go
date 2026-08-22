package orderhandler

import (
	"telegram-service-platform/delivery/telegramserver/messenger"
	"telegram-service-platform/service/checkoutservice"
	"telegram-service-platform/service/orderflowservice"
	"telegram-service-platform/service/pricingservice"
	"telegram-service-platform/service/productservice"
)

type Handler struct {
	productService   *productservice.Service
	checkoutService  *checkoutservice.Service
	orderFlowService *orderflowservice.Service
	pricingSvc       *pricingservice.Service
	messenger        messenger.Messenger
}

func New(
	productService *productservice.Service,
	checkoutService *checkoutservice.Service,
	orderFlowService *orderflowservice.Service,
	pricingSvc *pricingservice.Service,
	messenger messenger.Messenger,
) *Handler {
	return &Handler{
		productService:   productService,
		checkoutService:  checkoutService,
		orderFlowService: orderFlowService,
		messenger:        messenger,
		pricingSvc:       pricingSvc,
	}
}
