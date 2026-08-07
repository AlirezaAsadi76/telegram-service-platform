package app

import (
	"context"
	"log"
)

func (a *App) Start(ctx context.Context) error {

	go func() {

		if err := a.telegramBot.Start(ctx); err != nil {
			log.Println(err)
		}

	}()

	a.scheduler.Start()

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {

	if err := a.telegramBot.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.scheduler.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.redis.Close(); err != nil {
		return err
	}

	a.postgres.Close()

	return nil
}
