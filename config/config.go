package config

import (
	"time"
)

type HttpServer struct {
	Port int `koanf:"port"`
}

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
}

type Config struct {
	HttpServer  HttpServer  `koanf:"httpServer"`
	Application Application `koanf:"application"`
}
