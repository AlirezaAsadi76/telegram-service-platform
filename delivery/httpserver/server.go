package httpserver

import (
	"telegram-service-platform/config"

	"github.com/labstack/echo/v5"
)

type Server struct {
	Router *echo.Echo
	config config.Config
}

func New(cfg config.Config) *Server {
	e := echo.New()

	server := &Server{
		Router: e,
		config: cfg,
	}

	return server
}
