package orderfulfillmentrecoveryjob

import (
	"sync"
	"telegram-service-platform/service/orderfulfillmentservice"
	"telegram-service-platform/service/orderservice"
)

type Job struct {
	orderService        *orderservice.Service
	fulfillmentEnqueuer *orderfulfillmentservice.Service
	mutex               sync.Mutex
	config              Config
}

func New(
	orderService *orderservice.Service,
	fulfillmentEnqueuer *orderfulfillmentservice.Service,
	config Config,
) *Job {
	return &Job{
		orderService:        orderService,
		fulfillmentEnqueuer: fulfillmentEnqueuer,
		config:              config,
	}
}

func (j *Job) Name() string {
	return "order-fulfillment-recovery"
}
