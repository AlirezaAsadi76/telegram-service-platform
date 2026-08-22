package app

import (
	"telegram-service-platform/adapter/botadapter"
	"telegram-service-platform/adapter/exchangerate"
	"telegram-service-platform/adapter/fzrcards"
	"telegram-service-platform/adapter/fzrcards/telegramproduct"
	"telegram-service-platform/adapter/redisadapter"
	"telegram-service-platform/adapter/smm/justanotherpanel"
	"telegram-service-platform/config"
	"telegram-service-platform/delivery/telegramserver/messenger"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgresnotification"
	"telegram-service-platform/repository/postgresorder"
	"telegram-service-platform/repository/postgrespayment"
	"telegram-service-platform/repository/postgresproduct"
	"telegram-service-platform/repository/postgresprovider"
	"telegram-service-platform/repository/postgresuser"
	"telegram-service-platform/repository/postgreswallet"
	"telegram-service-platform/repository/redis/redisactivity"
	"telegram-service-platform/repository/redis/rediscatalog"
	"telegram-service-platform/repository/redis/redisidempotency"
	"telegram-service-platform/repository/redis/redisprice"
	"telegram-service-platform/repository/redis/redisqueue"
	"telegram-service-platform/repository/redis/redissmm"
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
	MessengerService    *messenger.Service
}

type Repositories struct {
	queueRepo redisqueue.DB
}

type Adapters struct {
	justPanelAdapter *justanotherpanel.Adapter
	botAdapter       *botadapter.Adapter
	postgresClient   *postgres.DB
	redisAdapter     *redisadapter.Adapter
}

func SetupDependencies(cfg config.Config) (*Dependencies, *Repositories, *Adapters) {

	// adapters
	justPanelAdapter := justanotherpanel.New(cfg.Justanotherpanel)
	botAdapter := botadapter.New(cfg.Telegram)
	postgresClient, nErr := postgres.New(cfg.Postgres)
	if nErr != nil {
		panic(nErr)
	}

	redisAdapter := redisadapter.New(cfg.RedisCli)

	// Repositories
	walletRepo := postgreswallet.New(postgresClient)
	paymentRepo := postgrespayment.New(postgresClient)
	orderRepo := postgresorder.New(postgresClient)
	providerRepo := postgresprovider.New(postgresClient)
	idempotencyRepo := redisidempotency.New(redisAdapter)
	queueRepo := redisqueue.New(redisAdapter)
	notificationRepo := postgresnotification.New(postgresClient)
	userRepo := postgresuser.New(postgresClient)
	priceRepo := redisprice.New(redisAdapter)
	productRepo := postgresproduct.New(postgresClient)
	catalogCache := rediscatalog.New(redisAdapter, cfg.CatalogCatch)
	smmCache := redissmm.New(redisAdapter, cfg.SmmRedis)
	activityTracker := redisactivity.New(redisAdapter, cfg.Activity)

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
	userSvc := userservice.New(walletSvc, userRepo, activityTracker)
	messengerService := messenger.New(botAdapter)

	productSvc := productservice.New(cfg.ProductService, pricingSvc, productRepo, catalogCache, smmCache, justPanelAdapter)
	// smmSvc.RegisterProvider("justanotherpanel", justanotherpanel.New(...))

	// Orchestrator
	// TODO: Replace nil messenger with actual implementation
	checkoutSvc := checkoutservice.New(walletSvc, paymentSvc, orderSvc, smmSvc, messengerService, idempotencyRepo, cfg.CheckoutSvc)

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
			MessengerService:    messengerService,
		},
		&Repositories{
			queueRepo: queueRepo,
		},
		&Adapters{
			botAdapter:       botAdapter,
			postgresClient:   postgresClient,
			redisAdapter:     redisAdapter,
			justPanelAdapter: justPanelAdapter,
		}
}
