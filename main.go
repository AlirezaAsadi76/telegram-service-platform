package main

import (
	"telegram-service-platform/config"
	"telegram-service-platform/repository/migrator"
)

func main() {
	cfg := config.Load("config.yml")

	mi := migrator.New(cfg.Postgres)

	if err := mi.Up(); err != nil {
		panic(err)
	}

}
