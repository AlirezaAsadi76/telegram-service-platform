package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"telegram-service-platform/app"
	"telegram-service-platform/config"
)

func main() {

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer cancel()
	cfg := config.Load("config.yml")
	fmt.Println("config : ", cfg)

	//mi := migrator.New(cfg.Postgres)
	//mi.Down()
	//if err := mi.Up(); err != nil {
	//	panic(err)
	//}

	application, err := app.New(cfg)

	if err != nil {
		panic(err)
	}

	err = application.Start(ctx)

	if err != nil {
		panic(err)
	}

	<-ctx.Done()

}
