package main

import (
	"context"
	"fmt"
	"telegram-service-platform/adapter/smm/justanotherpanel"
	"telegram-service-platform/config"
	"telegram-service-platform/params/smmprams"
)

func main() {
	cfg := config.Load("config.yml")
	fmt.Println("config : ", cfg)

	jpanelAdapter := justanotherpanel.New(cfg.Justanotherpanel)

	//if err, _ := jpanelAdapter.AllServices(context.Background()); err != nil {
	//	panic(err)
	//}
	ctx := context.Background()
	res, _ := jpanelAdapter.MultiStatus(ctx, []string{
		"1000091579",
	})

	fmt.Println(res)

	resr, errr := jpanelAdapter.MultiRefillStatus(ctx, []string{
		"996993167", "996987516", "996987338", "996986665", "996851554",
	})
	fmt.Println(errr)
	fmt.Println(resr)

	res1, gerr := jpanelAdapter.GetBalance(ctx)
	fmt.Println(gerr)
	fmt.Println(res1)
	res2, gerr2 := jpanelAdapter.Create(ctx, smmprams.CreateOrderAdapterRequest{
		ServiceID: "9079",
		Link:      "https://t.me/TelgramAdsPoplo",
		Quantity:  5000,
	})
	fmt.Println(gerr2)
	fmt.Println(res2)

	res4, gerr4 := jpanelAdapter.Cancel(ctx, []string{
		"996987338",
	})
	fmt.Println(gerr4)
	fmt.Println(res4)
}
