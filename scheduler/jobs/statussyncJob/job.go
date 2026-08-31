package statussyncjob

import (
	"sync"
	"telegram-service-platform/service/walletservice"

	"telegram-service-platform/service/notificationservice"
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/smmproviderservice"
)

type Job struct {
	orderService        *orderservice.Service
	smmProviderService  *smmproviderservice.Service
	notificationService *notificationservice.Service
	walletService       *walletservice.Service
	mutex               sync.Mutex
}

func New(
	orderService *orderservice.Service,
	smmProviderService *smmproviderservice.Service,
	notificationService *notificationservice.Service,
	walletService *walletservice.Service,

) *Job {
	return &Job{
		orderService:        orderService,
		smmProviderService:  smmProviderService,
		notificationService: notificationService,
		walletService:       walletService,
	}
}

func (j *Job) Name() string {
	return "status-sync"
}
