package producthandler

import "telegram-service-platform/delivery/telegramserver/messenger"

type Handler struct {
	productService ProductService
	messenger      messenger.Messenger
}

func New(productService ProductService, messenger messenger.Messenger) Handler {
	return Handler{
		productService: productService,
		messenger:      messenger,
	}
}
