package checkoutservice

import (
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/paymentservice"
	"telegram-service-platform/service/smmproviderservice"
	"telegram-service-platform/service/walletservice"
)

type Service struct {
	walletSvc           *walletservice.Service
	paymentSvc          *paymentservice.Service
	orderSvc            *orderservice.Service
	smmSvc              *smmproviderservice.Service
	walletPurchaseRepo  WalletPurchaseRepository
	fulfillmentEnqueuer FulfillmentEnqueuer
	messenger           Messenger
	idempotency         IdempotencyChecker
	config              Config
}

func New(
	walletSvc *walletservice.Service,
	paymentSvc *paymentservice.Service,
	orderSvc *orderservice.Service,
	smmSvc *smmproviderservice.Service,
	walletPurchaseRepo WalletPurchaseRepository,
	fulfillmentEnqueuer FulfillmentEnqueuer,
	messenger Messenger,
	idempotency IdempotencyChecker,
	config Config,
) *Service {
	return &Service{
		walletSvc:           walletSvc,
		paymentSvc:          paymentSvc,
		orderSvc:            orderSvc,
		smmSvc:              smmSvc,
		messenger:           messenger,
		idempotency:         idempotency,
		walletPurchaseRepo:  walletPurchaseRepo,
		fulfillmentEnqueuer: fulfillmentEnqueuer,
		config:              config,
	}
}
