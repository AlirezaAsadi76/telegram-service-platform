package main

import (
	"fmt"
	"telegram-service-platform/config"
	"telegram-service-platform/repository/postgres"
)

func main() {
	cfg := config.Load("config.yml")
	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("database connected")

}
