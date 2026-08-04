package userhandler

import "github.com/go-telegram/bot"

type Router struct {
	handler Handler
}

func NewRouter(
	handler Handler,
) Router {

	return Router{
		handler: handler,
	}
}
