package paymentexpiryjob

import (
	"sync"
	"telegram-service-platform/service/notificationservice"
	"telegram-service-platform/service/orderservice"
	"telegram-service-platform/service/paymentservice"
)

type Job struct {
	paymentService      paymentservice.Service
	orderService        orderservice.Service
	notificationService notificationservice.Service
	mutex               sync.Mutex
}

func New(
	paymentService paymentservice.Service,
	orderService orderservice.Service,
	notificationService notificationservice.Service,
) *Job {
	return &Job{
		paymentService:      paymentService,
		orderService:        orderService,
		notificationService: notificationService,
	}
}

func (j *Job) Name() string {
	return "payment-expiry"
}
