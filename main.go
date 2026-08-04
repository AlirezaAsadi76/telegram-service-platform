package main

import (
	"telegram-service-platform/config"
	"telegram-service-platform/delivery/telegramserver"
	"telegram-service-platform/delivery/telegramserver/handler/userhandler"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgresuser"
	"telegram-service-platform/service/userservice"
	"telegram-service-platform/validator/uservalidator"
)

func main() {
	cfg := config.Load("config.yml")

	//mi := migrator.New(cfg.Postgres)
	//
	//if err := mi.Up(); err != nil {
	//	panic(err)
	//}

	postgresRepo, nErr := postgres.New(cfg.Postgres)
	if nErr != nil {
		panic(nErr)
	}
	userRepo := postgresuser.New(postgresRepo)
	userSvc := userservice.New(userRepo)
	userValidator := uservalidator.New()
	userHandler := userhandler.New(userSvc, userValidator)
	bot, tErr := telegramserver.New(cfg.Telegram, userHandler)
	if tErr != nil {
		panic(tErr)
	}

}
