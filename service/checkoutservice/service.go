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
	transactionRepo     TransactionRepository
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
	transactionRepo TransactionRepository,
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
		transactionRepo:     transactionRepo,
		fulfillmentEnqueuer: fulfillmentEnqueuer,
		config:              config,
	}
}
