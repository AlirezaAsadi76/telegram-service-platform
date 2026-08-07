package config

import (
	"telegram-service-platform/adapter/exchangerate"
	"telegram-service-platform/adapter/fzrcards"
	"telegram-service-platform/adapter/redisadapter"
	"telegram-service-platform/delivery/telegramserver"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/scheduler"
	"telegram-service-platform/service/priceservice"
	"telegram-service-platform/service/productservice"
	"time"
)

type HttpServer struct {
	Port int `koanf:"port"`
}

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
}

type Config struct {
	HttpServer     HttpServer            `koanf:"httpServer"`
	Application    Application           `koanf:"application"`
	Postgres       postgres.DBConfig     `koanf:"postgres"`
	Telegram       telegramserver.Config `koanf:"telegram"`
	RedisCli       redisadapter.Config   `koanf:"redis"`
	ProductService productservice.Config `koanf:"productService"`
	PriceService   priceservice.Config   `koanf:"priceService"`
	Fzr            fzrcards.Config       `koanf:"fzr"`
	ExchangeRate   exchangerate.Config   `koanf:"exchangeRate"`
	Scheduler      scheduler.Config      `koanf:"scheduler"`
}
