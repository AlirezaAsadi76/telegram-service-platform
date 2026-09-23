package smmcatalogsyncjob

import (
	"sync"
)

type Job struct {
	productService ProductService
	mutex          sync.Mutex
}

func New(productService ProductService) *Job {
	return &Job{
		productService: productService,
	}
}

func (j *Job) Name() string {
	return "smm-catalog-sync"
}
