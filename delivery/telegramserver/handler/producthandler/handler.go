package producthandler

import "telegram-service-platform/service/productservice"

type Handler struct {
	productService productservice.Service
}

func New(productService productservice.Service) Handler {
	return Handler{
		productService: productService,
	}
}
