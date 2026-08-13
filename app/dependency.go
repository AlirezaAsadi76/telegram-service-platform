package app

import (
	"telegram-service-platform/adapter/redisadapter"
	"telegram-service-platform/config"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgresorder"
	"telegram-service-platform/repository/postgrespayment"
	"telegram-service-platform/repository/postgresprovider"
	"telegram-service-platform/repository/postgreswallet"
	"telegram-service-platform/service/checkoutservice"
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/paymentservice"
	"telegram-service-platform/service/smmproviderservice"
	"telegram-service-platform/service/walletservice"
)

type Dependencies struct {
	CheckoutService *checkoutservice.Service
	WalletService   *walletservice.Service
	OrderService    *orderservice.Service
	PaymentService  *paymentservice.Service
	SMMService      *smmproviderservice.Service
}

func SetupDependencies(cfg config.Config, pg *postgres.DB, redis *redisadapter.Redis) *Dependencies {
	// Repositories
	walletRepo := postgreswallet.New(pg)
	paymentRepo := postgrespayment.New(pg)
	orderRepo := postgresorder.New(pg)
	providerRepo := postgresprovider.New(pg)

	// Services
	walletSvc := walletservice.New(walletRepo, walletRepo, redis)
	paymentSvc := paymentservice.New(paymentRepo, nil, nil) // TODO: add adapters
	orderSvc := orderservice.New(orderRepo)
	smmSvc := smmproviderservice.New(providerRepo)
	// smmSvc.RegisterProvider("justanotherpanel", justanotherpanel.New(...))

	// Orchestrator
	// TODO: Replace nil messenger with actual implementation
	checkoutSvc := checkoutservice.New(walletSvc, paymentSvc, orderSvc, smmSvc, nil, redis)

	return &Dependencies{
		CheckoutService: checkoutSvc,
		WalletService:   walletSvc,
		OrderService:    orderSvc,
		PaymentService:  paymentSvc,
		SMMService:      smmSvc,
	}
}
