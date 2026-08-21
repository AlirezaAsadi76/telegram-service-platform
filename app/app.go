package app

import (
	"telegram-service-platform/adapter/botadapter"
	"telegram-service-platform/adapter/redisadapter"
	"telegram-service-platform/config"
	"telegram-service-platform/delivery/httpserver"
	"telegram-service-platform/delivery/telegramserver"
	"telegram-service-platform/delivery/telegramserver/handler/adminhandler"
	"telegram-service-platform/delivery/telegramserver/handler/mainhandler"
	"telegram-service-platform/delivery/telegramserver/handler/producthandler"
	"telegram-service-platform/delivery/telegramserver/handler/userhandler"
	"telegram-service-platform/delivery/telegramserver/messenger"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/scheduler"
	"telegram-service-platform/scheduler/jobs/notificationdispatchjob"
	"telegram-service-platform/scheduler/jobs/orderfulfillerjob"
	"telegram-service-platform/scheduler/jobs/paymentexpiryjob"
	"telegram-service-platform/scheduler/jobs/paymentverifyjob"
	"telegram-service-platform/scheduler/jobs/smmvalidationjob"
	statussyncjob "telegram-service-platform/scheduler/jobs/statussyncJob"
	"telegram-service-platform/validator/uservalidator"
)

type App struct {
	telegramBot   *telegramserver.Bot
	scheduler     *scheduler.Scheduler
	postgres      *postgres.DB
	redis         redisadapter.Adapter
	metricsServer *httpserver.Server
}

func New(cfg config.Config) (*App, error) {

	postgresClient, nErr := postgres.New(cfg.Postgres)
	if nErr != nil {
		panic(nErr)
	}

	redisAdapter := redisadapter.New(cfg.RedisCli)
	botAdapter := botadapter.New(cfg.Telegram)
	dependencies, repositories := SetupDependencies(cfg, postgresClient, &redisAdapter)

	userValidator := uservalidator.New()

	messengerService := messenger.New(botAdapter)
	// handler
	productHandler := producthandler.New(dependencies.ProductService, messengerService)
	userHandler := userhandler.New(dependencies.UserService, userValidator, messengerService)
	mainHandler := mainhandler.New(dependencies.ProductService, dependencies.UserService, messengerService)
	adminHandler := adminhandler.New(dependencies.CheckoutService, dependencies.UserService, messengerService, cfg.Admins)

	//priceRefreshJob := pricerefreshjob.New(dependencies.PriceService)
	pvj := paymentverifyjob.New(dependencies.PaymentService, dependencies.OrderService, dependencies.NotificationService, repositories.queueRepo, cfg.PaymentVerify)
	ofj := orderfulfillerjob.New(dependencies.OrderService, dependencies.SMMService, dependencies.NotificationService, repositories.queueRepo, cfg.OrderFulFiller)
	ssj := statussyncjob.New(dependencies.OrderService, dependencies.SMMService, dependencies.NotificationService)
	pej := paymentexpiryjob.New(dependencies.PaymentService, dependencies.OrderService, dependencies.NotificationService)
	ndj := notificationdispatchjob.New(dependencies.NotificationService, repositories.queueRepo, messengerService, cfg.NotificationJob)
	smj := smmvalidationjob.New(dependencies.ProductService, dependencies.NotificationService)
	schedulerObj, sErr := scheduler.New(cfg.Scheduler,
		pvj, ofj, ssj, pej, ndj, smj)
	if err := schedulerObj.Register(); err != nil {
		return nil, err
	}

	telegramBot, tErr := telegramserver.New(
		botAdapter,
		userHandler,
		productHandler,
		mainHandler,
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
			postgresClient.Connection().Ping, // یا متد ping مناسب روی postgres.DB
			redisAdapter.Ping,                // یا متد ping مناسب روی redis adapter
		)
	}

	return &App{

		telegramBot:   telegramBot,
		postgres:      postgresClient,
		redis:         redisAdapter,
		scheduler:     schedulerObj,
		metricsServer: metricsServer,
	}, nil

}
