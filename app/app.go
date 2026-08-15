package app

import (
	"telegram-service-platform/adapter/exchangerate"
	"telegram-service-platform/adapter/fzrcards"
	"telegram-service-platform/adapter/fzrcards/telegramproduct"
	"telegram-service-platform/adapter/redisadapter"
	"telegram-service-platform/config"
	"telegram-service-platform/delivery/httpserver"
	"telegram-service-platform/delivery/telegramserver"
	"telegram-service-platform/delivery/telegramserver/handler/producthandler"
	"telegram-service-platform/delivery/telegramserver/handler/userhandler"
	"telegram-service-platform/delivery/telegramserver/messenger"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgresproduct"
	"telegram-service-platform/repository/postgresuser"
	"telegram-service-platform/repository/redis/redisprice"
	"telegram-service-platform/scheduler"
	"telegram-service-platform/scheduler/jobs/notificationdispatchjob"
	"telegram-service-platform/scheduler/jobs/orderfulfillerjob"
	"telegram-service-platform/scheduler/jobs/paymentexpiryjob"
	"telegram-service-platform/scheduler/jobs/paymentverifyjob"
	"telegram-service-platform/scheduler/jobs/pricerefreshjob"
	statussyncjob "telegram-service-platform/scheduler/jobs/statussyncJob"
	"telegram-service-platform/service/priceservice"
	"telegram-service-platform/service/pricingservice"
	"telegram-service-platform/service/productservice"
	"telegram-service-platform/service/userservice"
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

	fzrClient := fzrcards.New(cfg.Fzr)
	telegramProvider := telegramproduct.New(fzrClient)
	exchangeRateProvider := exchangerate.New(cfg.ExchangeRate)

	userRepo := postgresuser.New(postgresClient)
	userSvc := userservice.New(userRepo)
	userValidator := uservalidator.New()

	priceRepo := redisprice.New(redisAdapter)

	priceService := priceservice.New(cfg.PriceService, priceRepo, telegramProvider, exchangeRateProvider)
	_ = priceService

	pricingSvc := pricingservice.New(priceRepo)

	productRepo := postgresproduct.New(postgresClient)
	productSvc := productservice.New(cfg.ProductService, pricingSvc, productRepo)

	messengerService := messenger.New()

	productHandler := producthandler.New(productSvc, messengerService)
	userHandler := userhandler.New(userSvc, userValidator, messengerService)

	priceRefreshJob := pricerefreshjob.New(priceService)
	pvj := paymentverifyjob.New(paymentService, orderService, notificationService, a.redis)
	ofj := orderfulfillerjob.New(a.orderService, a.smmProviderService, a.notificationService, a.redis)
	ssj := statussyncjob.New(a.orderService, a.smmProviderService, a.notificationService, a.redis)
	pej := paymentexpiryjob.New(a.paymentService, a.orderService, a.notificationService)
	ndj := notificationdispatchjob.New(a.notificationService, a.redis, a.telegramBotAdapter)

	schedulerObj, sErr := scheduler.New(cfg.Scheduler, priceRefreshJob)
	if err := schedulerObj.Register(); err != nil {
		return nil, err
	}

	telegramBot, tErr := telegramserver.New(
		cfg.Telegram,
		userHandler,
		productHandler,
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
			postgresClient.Pool.Ping, // یا متد ping مناسب روی postgres.DB
			redisAdapter.Ping,        // یا متد ping مناسب روی redis adapter
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
