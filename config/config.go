package config

import (
	"telegram-service-platform/repository/postgres"
	"time"
)

type HttpServer struct {
	Port int `koanf:"port"`
}

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
}

type Config struct {
	HttpServer  HttpServer        `koanf:"httpServer"`
	Application Application       `koanf:"application"`
	Postgres    postgres.DBConfig `koanf:"postgres"`
}
