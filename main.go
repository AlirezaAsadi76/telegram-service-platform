package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"telegram-service-platform/adapter/exchangerate"
	"telegram-service-platform/adapter/fzrcards"
	"telegram-service-platform/adapter/fzrcards/telegramproduct"
	"telegram-service-platform/adapter/redisadapter"
	"telegram-service-platform/config"
	"telegram-service-platform/delivery/telegramserver"
	"telegram-service-platform/delivery/telegramserver/handler/producthandler"
	"telegram-service-platform/delivery/telegramserver/handler/userhandler"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgresproduct"
	"telegram-service-platform/repository/postgresuser"
	"telegram-service-platform/repository/redis/redisprice"
	"telegram-service-platform/service/priceservice"
	"telegram-service-platform/service/pricingservice"
	"telegram-service-platform/service/productservice"
	"telegram-service-platform/service/userservice"
	"telegram-service-platform/validator/uservalidator"
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

	postgresRepo, nErr := postgres.New(cfg.Postgres)
	if nErr != nil {
		panic(nErr)
	}

	redisAdapter := redisadapter.New(cfg.RedisCli)

	fzrClient := fzrcards.New(cfg.Fzr)
	telegramProvider := telegramproduct.New(fzrClient)
	exchangeRateProvider := exchangerate.New(cfg.ExchangeRate)

	userRepo := postgresuser.New(postgresRepo)
	userSvc := userservice.New(userRepo)
	userValidator := uservalidator.New()
	userHandler := userhandler.New(userSvc, userValidator)

	priceRepo := redisprice.New(redisAdapter)

	priceService := priceservice.New(cfg.PriceService, priceRepo, telegramProvider, exchangeRateProvider)
	_ = priceService

	pricingSvc := pricingservice.New(priceRepo)

	productRepo := postgresproduct.New(postgresRepo)
	productSvc := productservice.New(cfg.ProductService, pricingSvc, productRepo)

	productHandler := producthandler.New(productSvc)

	telegramBot, tErr := telegramserver.New(
		cfg.Telegram,
		userHandler,
		productHandler,
	)

	if tErr != nil {
		panic(tErr)
	}

	telegramBot.Start(ctx)

}
