package app

import (
	"context"
	"log"
	"telegram-service-platform/logger"

	"go.uber.org/zap"
)

func (a *App) Start(ctx context.Context) error {

	if a.metricsServer != nil {
		if err := a.metricsServer.Start(ctx); err != nil {
			return err
		}
	}

	go func() {

		if err := a.telegramBot.Start(ctx); err != nil {
			log.Println(err)
		}

	}()

	a.scheduler.Start()

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {

	if a.metricsServer != nil {
		if err := a.metricsServer.Shutdown(ctx); err != nil {
			logger.Logger.Error("metrics server shutdown error", zap.Error(err))
		}
	}

	if err := a.telegramBot.Shutdown(ctx); err != nil {

		logger.Logger.Error("telegram server shutdown error", zap.Error(err))
		return err
	}

	if err := a.scheduler.Shutdown(ctx); err != nil {

		logger.Logger.Error("scheduler server shutdown error", zap.Error(err))
		return err
	}

	if err := a.redis.Close(); err != nil {
		logger.Logger.Error("redis server shutdown error", zap.Error(err))
		return err
	}

	a.postgres.Close()

	return nil
}
