package checkoutservice

import (
	"telegram-service-platform/service/notificationservice"
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
	notificationSvc     *notificationservice.Service
	walletPurchaseRepo  WalletPurchaseRepository
	fulfillmentEnqueuer FulfillmentEnqueuer
	idempotency         IdempotencyChecker
	config              Config
}

func New(
	walletSvc *walletservice.Service,
	paymentSvc *paymentservice.Service,
	orderSvc *orderservice.Service,
	smmSvc *smmproviderservice.Service,
	notificationSvc *notificationservice.Service,
	walletPurchaseRepo WalletPurchaseRepository,
	fulfillmentEnqueuer FulfillmentEnqueuer,
	idempotency IdempotencyChecker,
	config Config,
) *Service {
	return &Service{
		walletSvc:           walletSvc,
		paymentSvc:          paymentSvc,
		orderSvc:            orderSvc,
		notificationSvc:     notificationSvc,
		smmSvc:              smmSvc,
		idempotency:         idempotency,
		walletPurchaseRepo:  walletPurchaseRepo,
		fulfillmentEnqueuer: fulfillmentEnqueuer,
		config:              config,
	}
}
