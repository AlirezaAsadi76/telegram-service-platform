package main

import (
	"log"
	"telegram-service-platform/config"
	"telegram-service-platform/repository/migrator"
	"telegram-service-platform/repository/postgres"
)

func main() {
	cfg := config.Load("config.yml")

	postgresClient, err := postgres.New(cfg.PostgresTest)

	if err != nil {
		log.Fatal(err)
	}

	defer postgresClient.Close()

	mi := migrator.New(cfg.Postgres)
	mi.Down()
	if err := mi.Up(); err != nil {
		panic(err)
	}
}
