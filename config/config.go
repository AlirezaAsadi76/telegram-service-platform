package config

import (
	"telegram-service-platform/adapter/redisadapter"
	"telegram-service-platform/delivery/telegramserver"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/service/priceservice"
	"time"
)

type HttpServer struct {
	Port int `koanf:"port"`
}

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
}

type Config struct {
	HttpServer   HttpServer            `koanf:"httpServer"`
	Application  Application           `koanf:"application"`
	Postgres     postgres.DBConfig     `koanf:"postgres"`
	Telegram     telegramserver.Config `koanf:"telegram"`
	redisCli     redisadapter.Config   `koanf:"redis"`
	priceService priceservice.Config   `koanf:"priceService"`
}
