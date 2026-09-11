package main

import (
	"log"
	"telegram-service-platform/config"
	"telegram-service-platform/repository/migrator"
	"telegram-service-platform/repository/postgres"
)

func main() {
	cfg := config.Load("config.yml")

	cfgTestPostgres := postgres.DBConfig{
		Host:     cfg.PostgresTest.Host,
		Port:     cfg.PostgresTest.Port,
		User:     cfg.PostgresTest.User,
		Password: cfg.PostgresTest.Password,
		Database: cfg.PostgresTest.Database,
	}

	postgresClient, err := postgres.New(cfgTestPostgres)

	if err != nil {
		log.Fatal(err)
	}

	defer postgresClient.Close()

	mi := migrator.New(cfgTestPostgres)
	mi.Down()
	if err := mi.Up(); err != nil {
		panic(err)
	}
}
