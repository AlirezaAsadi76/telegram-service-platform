package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"telegram-service-platform/config"
	"telegram-service-platform/delivery/httpserver"
)

func main() {
	cfg := config.Load("config.yml")
	fmt.Println("config : ", cfg)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		cfg.Application.GracefulShutdownTimeout,
	)

	defer shutdownCancel()
	server := httpserver.New(cfg)
	if err := server.Start(shutdownCtx); err != nil {
		panic(err)
	}
	<-ctx.Done()
}
