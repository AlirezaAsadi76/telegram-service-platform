package main

import (
	"context"
	"fmt"
	"telegram-service-platform/adapter/smm/justanotherpanel"
	"telegram-service-platform/config"
)

func main() {
	cfg := config.Load("config.yml")
	fmt.Println("config : ", cfg)

	jpanelAdapter := justanotherpanel.New(cfg.Justanotherpanel)

	//if err, _ := jpanelAdapter.AllServices(context.Background()); err != nil {
	//	panic(err)
	//}
	ctx := context.Background()
	jpanelAdapter.Status(ctx)
}
