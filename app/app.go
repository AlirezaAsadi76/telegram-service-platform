package app

import (
	"telegram-service-platform/adapter/redisadapter"
	"telegram-service-platform/config"
	"telegram-service-platform/delivery/httpserver"
	"telegram-service-platform/delivery/telegramserver"
	"telegram-service-platform/delivery/telegramserver/dispatcher"
	"telegram-service-platform/delivery/telegramserver/handler/adminhandler"
	"telegram-service-platform/delivery/telegramserver/handler/mainhandler"
	"telegram-service-platform/delivery/telegramserver/handler/orderhandler"
	"telegram-service-platform/delivery/telegramserver/handler/producthandler"
	"telegram-service-platform/delivery/telegramserver/handler/userhandler"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/scheduler"
	"telegram-service-platform/scheduler/jobs/notificationdispatchjob"
	"telegram-service-platform/scheduler/jobs/orderfulfillerjob"
	"telegram-service-platform/scheduler/jobs/paymentexpiryjob"
	"telegram-service-platform/scheduler/jobs/paymentverifyjob"
	"telegram-service-platform/scheduler/jobs/pricerefreshjob"
	"telegram-service-platform/scheduler/jobs/smmvalidationjob"
	statussyncjob "telegram-service-platform/scheduler/jobs/statussyncJob"
	"telegram-service-platform/validator/ordervalidator"
	"telegram-service-platform/validator/uservalidator"
)

type App struct {
	telegramBot   *telegramserver.Bot
	scheduler     *scheduler.Scheduler
	postgres      *postgres.DB
	redis         *redisadapter.Adapter
	metricsServer *httpserver.Server
}

func New(cfg config.Config) (*App, error) {

	dependencies, repositories, adapters := SetupDependencies(cfg)

	userValidator := uservalidator.New()
	orderValidator := ordervalidator.New()

	// handler
	productHandler := producthandler.New(dependencies.ProductService, dependencies.MessengerService)
	userHandler := userhandler.New(dependencies.UserService, userValidator, dependencies.MessengerService)
	mainHandler := mainhandler.New(dependencies.ProductService, dependencies.UserService, dependencies.OrderFlowService, dependencies.MessengerService)
	adminHandler := adminhandler.New(dependencies.CheckoutService, dependencies.UserService, dependencies.MessengerService, cfg.Admins)
	orderHandler := orderhandler.New(
		dependencies.ProductService, dependencies.CheckoutService,
		dependencies.OrderFlowService, dependencies.PricingService,
		dependencies.UserService, dependencies.MessengerService, orderValidator)

	//conversations
	conversationDispatcher := dispatcher.New(
		orderHandler,
	)

	prj := pricerefreshjob.New(dependencies.PriceService)
	pvj := paymentverifyjob.New(dependencies.PaymentService, dependencies.OrderService, dependencies.NotificationService, repositories.queueRepo, cfg.PaymentVerify)
	ofj := orderfulfillerjob.New(dependencies.OrderService, dependencies.SMMService, dependencies.NotificationService, dependencies.WalletService, repositories.queueRepo, cfg.OrderFulFiller)
	ssj := statussyncjob.New(dependencies.OrderService, dependencies.SMMService, dependencies.NotificationService, dependencies.WalletService)
	pej := paymentexpiryjob.New(dependencies.PaymentService, dependencies.OrderService, dependencies.NotificationService)
	ndj := notificationdispatchjob.New(dependencies.NotificationService, repositories.queueRepo, dependencies.MessengerService, cfg.NotificationJob)
	smj := smmvalidationjob.New(dependencies.ProductService, dependencies.NotificationService)
	schedulerObj, sErr := scheduler.New(cfg.Scheduler,
		prj, pvj, ofj, ssj, pej, ndj, smj)
	if err := schedulerObj.Register(); err != nil {
		return nil, err
	}

	telegramBot, tErr := telegramserver.New(
		adapters.botAdapter,
		conversationDispatcher,
		orderHandler,
		mainHandler,
		userHandler,
		productHandler,
		adminHandler,
	)

	if tErr != nil {
		panic(tErr)
	}
	if sErr != nil {
		panic(sErr)
	}

	var metricsServer *httpserver.Server
	if cfg.MetricsServer.Enabled {
		metricsServer = httpserver.New(
			cfg.MetricsServer.Port,
			adapters.postgresClient.Connection().Ping, // یا متد ping مناسب روی postgres.DB
			adapters.redisAdapter.Ping,                // یا متد ping مناسب روی redis adapter
		)
	}

	return &App{

		telegramBot:   telegramBot,
		postgres:      adapters.postgresClient,
		redis:         adapters.redisAdapter,
		scheduler:     schedulerObj,
		metricsServer: metricsServer,
	}, nil

}
