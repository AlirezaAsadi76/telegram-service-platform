package app

import "context"

func (a *App) Start(ctx context.Context) error {

	go func() {

		a.telegram.Start(ctx)

	}()

	a.scheduler.Start(ctx)

	return nil
}
