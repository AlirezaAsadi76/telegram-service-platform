package main

import (
	"context"
	"log"

	"telegram-service-platform/config"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/seeder/product"
)

func main() {

	cfg := config.Load("config.yml")

	ctx := context.Background()

	db, err := postgres.New(cfg.Postgres)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := product.SeedStarPlans(
		ctx,
		db.Connection(),
	); err != nil {
		log.Fatal(err)
	}

	if err := product.SeedPremiumPlans(
		ctx,
		db.Connection(),
	); err != nil {
		log.Fatal(err)
	}

	log.Println("seed completed")
}
