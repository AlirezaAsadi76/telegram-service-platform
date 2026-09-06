package main

import (
	"context"
	"log"
	"telegram-service-platform/config"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgresorder"
	"telegram-service-platform/repository/postgresuser"
	"telegram-service-platform/repository/postgreswallet"
	seeder2 "telegram-service-platform/repository/seeder"
)

func main() {

	cfg := config.Load("config.yml")

	ctx := context.Background()

	postgresClient, err := postgres.New(cfg.Postgres)

	if err != nil {
		log.Fatal(err)
	}

	defer postgresClient.Close()

	//mi := migrator.New(cfg.Postgres)
	//mi.Down()
	//if err := mi.Up(); err != nil {
	//	panic(err)
	//}
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

	//if err := smmseeder.SeedSMMData(ctx, db); err != nil {
	//	log.Fatal(err)
	//}
	walletRepo := postgreswallet.New(postgresClient)
	orderRepo := postgresorder.New(postgresClient)
	userRepo := postgresuser.New(postgresClient)

	seeder := seeder2.New(userRepo, walletRepo, orderRepo)

	seeder.SeedTestUser(ctx)
	seeder.SeedTestWallet(ctx)
	seeder.SeedTestOrders(ctx)

	log.Println("seed completed")
}
