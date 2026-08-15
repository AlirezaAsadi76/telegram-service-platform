package app

import (
	"telegram-service-platform/adapter/exchangerate"
	"telegram-service-platform/adapter/fzrcards"
	"telegram-service-platform/adapter/fzrcards/telegramproduct"
	"telegram-service-platform/adapter/redisadapter"
	"telegram-service-platform/config"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgresnotification"
	"telegram-service-platform/repository/postgresorder"
	"telegram-service-platform/repository/postgrespayment"
	"telegram-service-platform/repository/postgresproduct"
	"telegram-service-platform/repository/postgresprovider"
	"telegram-service-platform/repository/postgresuser"
	"telegram-service-platform/repository/postgreswallet"
	"telegram-service-platform/repository/redis/redisidempotency"
	"telegram-service-platform/repository/redis/redisprice"
	"telegram-service-platform/repository/redis/redisqueue"
	"telegram-service-platform/service/checkoutservice"
	"telegram-service-platform/service/notificationservice"
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/paymentservice"
	"telegram-service-platform/service/priceservice"
	"telegram-service-platform/service/pricingservice"
	"telegram-service-platform/service/productservice"
	"telegram-service-platform/service/smmproviderservice"
	"telegram-service-platform/service/userservice"
	"telegram-service-platform/service/walletservice"
)

type Dependencies struct {
	CheckoutService     *checkoutservice.Service
	WalletService       *walletservice.Service
	OrderService        *orderservice.Service
	PaymentService      *paymentservice.Service
	SMMService          *smmproviderservice.Service
	NotificationService *notificationservice.Service
	UserService         *userservice.Service
	PricingService      *pricingservice.Service
	PriceService        *priceservice.Service
	ProductService      *productservice.Service
}

type Repositories struct {
	queueRepo redisqueue.DB
}

func SetupDependencies(cfg config.Config, pg *postgres.DB, redis *redisadapter.Adapter) (*Dependencies, *Repositories) {
	// Repositories
	walletRepo := postgreswallet.New(pg)
	paymentRepo := postgrespayment.New(pg)
	orderRepo := postgresorder.New(pg)
	providerRepo := postgresprovider.New(pg)
	idempotencyRepo := redisidempotency.New(redis)
	queueRepo := redisqueue.New(redis)
	notificationRepo := postgresnotification.New(pg)
	userRepo := postgresuser.New(pg)
	priceRepo := redisprice.New(redis)
	productRepo := postgresproduct.New(pg)
	// Providers
	fzrClient := fzrcards.New(cfg.Fzr)
	telegramProvider := telegramproduct.New(fzrClient)
	exchangeRateProvider := exchangerate.New(cfg.ExchangeRate)

	// Services
	walletSvc := walletservice.New(walletRepo, walletRepo, idempotencyRepo, cfg.WalletSvc)
	paymentSvc := paymentservice.New(paymentRepo, nil, nil) // TODO: add adapters
	orderSvc := orderservice.New(orderRepo)
	smmSvc := smmproviderservice.New(providerRepo, cfg.SmmSvc)
	notificationSVC := notificationservice.New(notificationRepo, queueRepo, cfg.NotificationSvc)
	priceService := priceservice.New(cfg.PriceService, priceRepo, telegramProvider, exchangeRateProvider)
	pricingSvc := pricingservice.New(priceRepo)
	userSvc := userservice.New(userRepo)
	productSvc := productservice.New(cfg.ProductService, pricingSvc, productRepo)
	// smmSvc.RegisterProvider("justanotherpanel", justanotherpanel.New(...))

	// Orchestrator
	// TODO: Replace nil messenger with actual implementation
	checkoutSvc := checkoutservice.New(walletSvc, paymentSvc, orderSvc, smmSvc, nil, idempotencyRepo, cfg.CheckoutSvc)

	return &Dependencies{
			CheckoutService:     checkoutSvc,
			WalletService:       walletSvc,
			OrderService:        orderSvc,
			PaymentService:      paymentSvc,
			SMMService:          smmSvc,
			NotificationService: notificationSVC,
			UserService:         userSvc,
			PricingService:      pricingSvc,
			PriceService:        priceService,
			ProductService:      productSvc,
		},
		&Repositories{
			queueRepo: queueRepo,
		}
}
