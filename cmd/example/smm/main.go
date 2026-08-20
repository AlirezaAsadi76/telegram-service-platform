package main

import (
	"context"
	"fmt"
	"telegram-service-platform/config"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgresproduct"
	"telegram-service-platform/service/productservice"
)

func main() {
	cfg := config.Load("config.yml")
	fmt.Println("config : ", cfg)
	ctx := context.Background()
	//jpanelAdapter := justanotherpanel.New(cfg.Justanotherpanel)
	//
	//ff, err := jpanelAdapter.AllServices(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(ff.Services)
	//ctx := context.Background()
	//res, _ := jpanelAdapter.MultiStatus(ctx, []string{
	//	"1000091579",
	//})
	//
	//fmt.Println(res)
	//
	//resr, errr := jpanelAdapter.MultiRefillStatus(ctx, []string{
	//	"996993167", "996987516", "996987338", "996986665", "996851554",
	//})
	//fmt.Println(errr)
	//fmt.Println(resr)
	//
	//res1, gerr := jpanelAdapter.GetBalance(ctx)
	//fmt.Println(gerr)
	//fmt.Println(res1)
	//res2, gerr2 := jpanelAdapter.Create(ctx, smmprams.CreateOrderAdapterRequest{
	//	ServiceID: "34134234324342341",
	//	Link:      "https://t.me/TelgramAdsPoplo",
	//	Quantity:  5000,
	//})
	//fmt.Println(gerr2)
	//fmt.Println(res2)
	//
	//res4, gerr4 := jpanelAdapter.Cancel(ctx, []string{
	//	"996987338",
	//})
	//fmt.Println(gerr4)
	//fmt.Println(res4)

	//mi := migrator.New(cfg.Postgres)
	//mi.Down()
	//if err := mi.Up(); err != nil {
	//	panic(err)
	//}

	pq, _ := postgres.New(cfg.Postgres)

	repoProduct := postgresproduct.New(pq)

	productSvc := productservice.New(cfg.ProductService, nil, repoProduct)

	//{
	//	"service": 7102,
	//	"name": "Telegram Members [Refill: 7D] [Start Time: 0 - 1 Hr] [Max: 30K] [Speed: 20K/D] 💧⛔",
	//	"type": "Default",
	//	"rate": "0.125",
	//	"min": 10,
	//	"max": 100000,
	//	"dripfeed": true,
	//	"refill": false,
	//	"cancel": true,
	//	"category": "Telegram Members"
	//},

	//
	//{
	//	"service": 5994,
	//	"name": "Instagram Views [Max: 10M] [Start Time: 0-1 Hour] [Speed: 200K/D]",
	//	"type": "Default",
	//	"rate": "0.0015",
	//	"min": 100,
	//	"max": 2147483647,
	//	"dripfeed": true,
	//	"refill": false,
	//	"cancel": false,
	//	"category": "Instagram Views"
	//}
	//iErr := repoProduct.SMMServiceCreateOrUpdate(ctx, smmentity.SMM{
	//	Service:      5994,
	//	Name:         "Instagram Views [Max: 10M] [Start Time: 0-1 Hour] [Speed: 200K/D]",
	//	Type:         "Default",
	//	Rate:         "0.0015",
	//	Min:          100,
	//	Max:          1000000,
	//	DripFeed:     true,
	//	Refill:       false,
	//	Cancel:       false,
	//	IsActive:     true,
	//	CategoryType:     "Instagram Views",
	//	ProviderName: "justanotherpanel",
	//})
	//cErr := repoProduct.SMMServiceCreateOrUpdate(ctx, smmentity.SMM{
	//	Service:      7102,
	//	Name:         "Telegram Members [Refill: 7D] [Start Time: 0 - 1 Hr] [Max: 30K] [Speed: 20K/D] 💧⛔",
	//	Type:         "Default",
	//	Rate:         "0.125",
	//	Min:          10,
	//	Max:          100000,
	//	DripFeed:     true,
	//	Refill:       false,
	//	Cancel:       true,
	//	IsActive:     true,
	//	CategoryType:     "Telegram Members",
	//	ProviderName: "justanotherpanel",
	//})
	//tErr := repoProduct.SMMServiceCreateOrUpdate(ctx, smmentity.SMM{
	//	Service:      10338,
	//	Name:         "TikTok Followers [Refill: No] [Max: 100K] [Start Time: 0 - 2 Hours] [Speed: 20K/Day] ⛔️💧",
	//	Type:         "Default",
	//	Rate:         "1.17",
	//	Min:          100,
	//	Max:          100000,
	//	DripFeed:     true,
	//	Refill:       false,
	//	Cancel:       false,
	//	IsActive:     true,
	//	CategoryType:     "Tiktok Followers",
	//	ProviderName: "justanotherpanel",
	//})
	//fmt.Println(tErr, iErr, cErr)
	//
	//{
	//	"service": 10338,
	//	"name": "TikTok Followers [Refill: No] [Max: 100K] [Start Time: 0 - 2 Hours] [Speed: 20K/Day] ⛔️💧",
	//	"type": "Default",
	//	"rate": "1.17",
	//	"min": 10,
	//	"max": 100000,
	//	"dripfeed": true,
	//	"refill": false,
	//	"cancel": true,
	//	"category": "Tiktok Followers"
	//},

	//s, e := repoProduct.SMMServiceGetByD(ctx, 1)
	//sms, err := repoProduct.SMMServiceGetAll(ctx)
	//fmt.Println(s, e)
	//fmt.Println(sms, err)

	//productSvc.CreateSMMMapping(ctx, productparams.CreateSMMMappingRequest{
	//	SmmServiceId: 1,
	//	Name:         "Instagram Views",
	//	Platform:     smmentity.InstagramPlatform,
	//	CategoryType:     smmentity.ViewCategory,
	//	Description:  "همینطوری",
	//	IsActive:     true,
	//	ButtonName:   "عادی",
	//})
	//productSvc.CreateSMMMapping(ctx, productparams.CreateSMMMappingRequest{
	//	SmmServiceId: 2,
	//	Name:         "Telegram Members",
	//	Platform:     smmentity.TelegramPlatform,
	//	CategoryType:     smmentity.MemberCategory,
	//	Description:  "همینطوری",
	//	IsActive:     true,
	//	ButtonName:   "عادی",
	//})
	//productSvc.CreateSMMMapping(ctx, productparams.CreateSMMMappingRequest{
	//	SmmServiceId: 3,
	//	Name:         "TikTok Followers",
	//	Platform:     smmentity.TickTockPlatform,
	//	CategoryType:     smmentity.FollowerCategory,
	//	Description:  "همینطوری",
	//	IsActive:     true,
	//	ButtonName:   "عادی",
	//})

	catalog, gcErr := productSvc.GetSMMCatalog(ctx)
	fmt.Println(catalog, gcErr)
	smm, giErr := productSvc.GetSMMMappingByID(ctx, productparams.GetSmmMappingByIDRequest{
		Id: 1,
	})
	fmt.Println(smm, giErr)
	plat, pcErr := productSvc.GetSMMMappingsByPlatformCategory(ctx, productparams.GetSmmMappingByPlatformCategoryRequest{
		Platform: smmentity.TickTockPlatform,
		Category: smmentity.FollowerCategory,
	})
	fmt.Println(plat, pcErr)
	platr, pcErrr := productSvc.GetSMMMappingsByPlatformCategory(ctx, productparams.GetSmmMappingByPlatformCategoryRequest{
		Platform: smmentity.TickTockPlatform,
		Category: smmentity.ViewCategory,
	})
	fmt.Println(platr, pcErrr)

}
