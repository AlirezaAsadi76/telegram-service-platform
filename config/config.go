package config

import (
	"telegram-service-platform/adapter/exchangerate"
	"telegram-service-platform/adapter/fzrcards"
	"telegram-service-platform/adapter/redisadapter"
	"telegram-service-platform/adapter/smm/justanotherpanel"
	"telegram-service-platform/delivery/telegramserver"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/scheduler"
	"telegram-service-platform/scheduler/jobs/notificationdispatchjob"
	"telegram-service-platform/scheduler/jobs/orderfulfillerjob"
	"telegram-service-platform/scheduler/jobs/paymentverifyjob"
	"telegram-service-platform/service/checkoutservice"
	"telegram-service-platform/service/notificationservice"
	"telegram-service-platform/service/priceservice"
	"telegram-service-platform/service/productservice"
	"telegram-service-platform/service/smmproviderservice"
	"telegram-service-platform/service/walletservice"
	"time"
)

type HttpServer struct {
	Port int `koanf:"port"`
}

type MetricsServer struct {
	Port    int  `koanf:"port"`
	Enabled bool `koanf:"enabled"`
}

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
}

type Config struct {
	HttpServer       HttpServer                     `koanf:"httpServer"`
	MetricsServer    MetricsServer                  `koanf:"metricsServer"`
	Application      Application                    `koanf:"application"`
	Postgres         postgres.DBConfig              `koanf:"postgres"`
	Telegram         telegramserver.Config          `koanf:"telegram"`
	RedisCli         redisadapter.Config            `koanf:"redis"`
	ProductService   productservice.Config          `koanf:"productService"`
	PriceService     priceservice.Config            `koanf:"priceService"`
	Fzr              fzrcards.Config                `koanf:"fzr"`
	ExchangeRate     exchangerate.Config            `koanf:"exchangeRate"`
	Scheduler        scheduler.Config               `koanf:"scheduler"`
	WalletSvc        walletservice.Config           `koanf:"wallet"`
	SmmSvc           smmproviderservice.Config      `koanf:"smmProviderSvc"`
	CheckoutSvc      checkoutservice.Config         `koanf:"checkoutSvc"`
	NotificationSvc  notificationservice.Config     `koanf:"notificationSvc"`
	PaymentVerify    paymentverifyjob.Config        `koanf:"paymentVerify"`
	OrderFulFiller   orderfulfillerjob.Config       `koanf:"orderFulFiller"`
	NotificationJob  notificationdispatchjob.Config `koanf:"notificationJob"`
	Justanotherpanel justanotherpanel.Config        `koanf:"justanotherPanel"`
}
