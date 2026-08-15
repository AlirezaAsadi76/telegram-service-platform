package paymentverifyjob

import (
	"sync"

	"telegram-service-platform/service/notificationservice"
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/paymentservice"
)

type Job struct {
	paymentService      *paymentservice.Service
	orderService        *orderservice.Service
	notificationService *notificationservice.Service
	redis               RedisRepository
	mutex               sync.Mutex
	config              Config
}

func New(
	paymentService *paymentservice.Service,
	orderService *orderservice.Service,
	notificationService *notificationservice.Service,
	redis RedisRepository,
	config Config,
) *Job {
	return &Job{
		paymentService:      paymentService,
		orderService:        orderService,
		notificationService: notificationService,
		redis:               redis,
		config:              config,
	}
}

func (j *Job) Name() string {
	return "payment-verify"
}
