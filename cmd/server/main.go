package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"telegram-service-platform/app"
	"telegram-service-platform/config"
	"telegram-service-platform/delivery/httpserver"
	"telegram-service-platform/delivery/httpserver/userhandler"
	"telegram-service-platform/logger"
	"telegram-service-platform/validator/uservalidator"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load("config.yml")
	fmt.Println("config loaded successfully")

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	dependencies, _, _ := app.SetupDependencies(cfg)
	userVal := uservalidator.New(dependencies.UserService)

	userHandler := userhandler.New(
		dependencies.UserService,
		dependencies.AuthService,
		dependencies.WalletService,
		userVal,
		cfg.Auth,
		cfg.Telegram.Token,
	)

	server := httpserver.New(cfg, userHandler)

	logger.Logger.Info("starting application...")

	if err := server.Start(ctx); err != nil {
		if err != http.ErrServerClosed {
			logger.Logger.Fatal("failed to start HTTP server", zap.Error(err))
		}
	}

	logger.Logger.Info("server stopped gracefully")
}
