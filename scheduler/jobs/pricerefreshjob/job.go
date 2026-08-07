package pricerefreshjob

import (
	"sync"
	"telegram-service-platform/service/priceservice"
)

type Job struct {
	priceService priceservice.Service
	mutex        sync.Mutex
}

func New(priceService priceservice.Service) Job {
	return Job{
		priceService: priceService,
	}
}

func (j *Job) Name() string {
	return "price-refresh"
}
