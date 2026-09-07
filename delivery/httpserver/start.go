package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"telegram-service-platform/logger"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
)

func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.config.HttpServer.Port)
	logger.Logger.Info("starting HTTP server", zap.String("address", addr))

	s.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			logger.Logger.Info("HTTP request",
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))
	s.Router.Use(middleware.Recover())

	s.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"https://web.telegram.org",
			"http://localhost:5173",
			"http://localhost:3000",
		},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Telegram-Init-Data"},
		ExposeHeaders:    []string{echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	for _, handler := range s.handlers {
		handler.SetRoutes(s.Router)
	}

	startConfig := echo.StartConfig{
		Address:         fmt.Sprintf(":%d", s.config.HttpServer.Port),
		GracefulTimeout: s.config.Application.GracefulShutdownTimeout,
	}

	if err := startConfig.Start(ctx, s.Router); err != nil {
		s.Router.Logger.Error("failed to start server", "error", err)
		return err
	}
	return nil
}
