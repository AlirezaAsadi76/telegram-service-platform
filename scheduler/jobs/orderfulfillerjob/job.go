package orderfulfillerjob

import (
	"sync"
	"telegram-service-platform/service/notificationservice"
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/smmproviderservice"
)

type Job struct {
	orderService        *orderservice.Service
	smmProviderService  *smmproviderservice.Service
	notificationService *notificationservice.Service
	redis               RedisRepository
	mutex               sync.Mutex
	config              Config
}

func New(
	orderService *orderservice.Service,
	smmProviderService *smmproviderservice.Service,
	notificationService *notificationservice.Service,
	redis RedisRepository,
	config Config,
) *Job {
	return &Job{
		orderService:        orderService,
		smmProviderService:  smmProviderService,
		notificationService: notificationService,
		redis:               redis,
		config:              config,
	}
}

func (j *Job) Name() string {
	return "order-fulfiller"
}
