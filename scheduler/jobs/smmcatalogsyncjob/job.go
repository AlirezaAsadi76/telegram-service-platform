package smmcatalogsyncjob

import (
	"sync"

	"telegram-service-platform/service/productservice"
)

type Job struct {
	productService ProductService
	mutex          sync.Mutex
}

func New(productService *productservice.Service) *Job {
	return &Job{
		productService: productService,
	}
}

func (j *Job) Name() string {
	return "smm-catalog-sync"
}
