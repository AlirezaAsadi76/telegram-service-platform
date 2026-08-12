package app

import (
	"telegram-service-platform/adapter/redisadapter"
	"telegram-service-platform/config"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgresorder"
	"telegram-service-platform/repository/postgrespayment"
	"telegram-service-platform/repository/postgresprovider"
	"telegram-service-platform/repository/postgreswallet"
	"telegram-service-platform/repository/redis/redisidempotency"
	"telegram-service-platform/service/checkoutservice"
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/paymentservice"
	"telegram-service-platform/service/smmproviderservice"
	"telegram-service-platform/service/walletservice"
)

type Dependencies struct {
	CheckoutService *checkoutservice.Service
}

func SetupDependencies(cfg config.Config, pg *postgres.DB, redis *redisadapter.Adapter) *Dependencies {
	// Repositories

	idempotency := redisidempotency.New(redis)

	walletRepo := postgreswallet.New(pg)
	paymentRepo := postgrespayment.New(pg)
	orderRepo := postgresorder.New(pg)
	providerRepo := postgresprovider.New(pg)

	// Services
	walletSvc := walletservice.New(walletRepo, walletRepo, idempotency)
	paymentSvc := paymentservice.New(paymentRepo, zarinpalAdapter, cryptoAdapter)
	orderSvc := orderservice.New(orderRepo)
	smmSvc := smmproviderservice.New(providerRepo)
	// smmSvc.RegisterProvider(...)

	// Orchestrator
	checkoutSvc := checkoutservice.New(walletSvc, paymentSvc, orderSvc, smmSvc, messenger, redis)

	return &Dependencies{CheckoutService: checkoutSvc}
}
