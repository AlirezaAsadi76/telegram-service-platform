package main

import (
	"context"
	"log"
	"telegram-service-platform/repository/migrator"
	"telegram-service-platform/repository/seeder/smmseeder"

	"telegram-service-platform/config"
	"telegram-service-platform/repository/postgres"
)

func main() {

	cfg := config.Load("config.yml")

	ctx := context.Background()

	db, err := postgres.New(cfg.Postgres)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	mi := migrator.New(cfg.Postgres)
	mi.Down()
	if err := mi.Up(); err != nil {
		panic(err)
	}
	//if err := product.SeedStarPlans(
	//	ctx,
	//	db.Connection(),
	//); err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err := product.SeedPremiumPlans(
	//	ctx,
	//	db.Connection(),
	//); err != nil {
	//	log.Fatal(err)
	//}

	if err := smmseeder.SeedSMMData(ctx, db); err != nil {
		log.Fatal(err)
	}

	log.Println("seed completed")
}
