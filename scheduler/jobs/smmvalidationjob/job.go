package smmvalidationjob

import (
	"context"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/params/notificationparams"
)

type NotificationService interface {
	Create(ctx context.Context, req notificationparams.CreateRequest) error
}

type ProductService interface {
	GetMissingSMMServices(ctx context.Context) ([]smmentity.SMM, error)
}

type SMMAdapterInterface interface {
	AllServices(ctx context.Context) ([]smmentity.SMM, error)
}

// TODO- later read admins id as database
type Config struct {
	AdminUserID int64 // شناسه کاربری ادمین برای ارسال هشدار
}

type Job struct {
	productService  ProductService
	notificationSvc NotificationService
	config          Config
}

func New(productSvc ProductService, notificationSvc NotificationService) *Job {
	cfg := Config{
		AdminUserID: 962404032,
	}
	return &Job{
		productService:  productSvc,
		notificationSvc: notificationSvc,
		config:          cfg,
	}
}

func (j *Job) Name() string {
	return "smm-validation"
}
