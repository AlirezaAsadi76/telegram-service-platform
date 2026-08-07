package pricerefreshjob

import "telegram-service-platform/service/priceservice"

type Job struct {
	priceService priceservice.Service
}

func New(priceService priceservice.Service) Job {
	return Job{
		priceService: priceService,
	}
}

func (j Job) Name() string {
	return "price-refresh"
}
