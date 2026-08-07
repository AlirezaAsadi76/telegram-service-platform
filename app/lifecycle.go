package app

import "context"

func (a *App) Start(ctx context.Context) error {

	go func() {

		a.telegramBot.Start(ctx)

	}()

	if err := a.scheduler.Shutdown(); err != nil {
		return err
	}

	_ = a.redis.Close()

	a.postgres.Close()

	return nil

}
