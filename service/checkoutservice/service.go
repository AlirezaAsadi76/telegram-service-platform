package checkoutservice

import (
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/paymentservice"
	"telegram-service-platform/service/smmproviderservice"
	"telegram-service-platform/service/walletservice"
)

type Service struct {
	walletSvc  *walletservice.Service
	paymentSvc *paymentservice.Service
	orderSvc   *orderservice.Service
	smmSvc     *smmproviderservice.Service
	messenger  Messenger
}

func New(
	walletSvc *walletservice.Service,
	paymentSvc *paymentservice.Service,
	orderSvc *orderservice.Service,
	smmSvc *smmproviderservice.Service,
	messenger Messenger,
) *Service {
	return &Service{
		walletSvc:  walletSvc,
		paymentSvc: paymentSvc,
		orderSvc:   orderSvc,
		smmSvc:     smmSvc,
		messenger:  messenger,
	}
}
