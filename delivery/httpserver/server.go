package httpserver

import (
	"telegram-service-platform/config"

	"github.com/labstack/echo/v5"
)

type Server struct {
	Router   *echo.Echo
	handlers []Handlers
	config   config.Config
}

func New(cfg config.Config, handlers ...Handlers) *Server {
	e := echo.New()

	server := &Server{
		Router:   e,
		handlers: handlers,
		config:   cfg,
	}

	return server
}
