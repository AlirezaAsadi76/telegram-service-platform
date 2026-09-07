package userhandler

import (
	"telegram-service-platform/delivery/httpserver/middleware"

	"github.com/labstack/echo/v5"
)

func (h *Handler) SetRoutes(e *echo.Echo) {
	userGroup := e.Group("/user")
	userGroup.POST("/register", h.registerHandler)
	userGroup.POST("/login", h.loginHandler)
	userGroup.GET("/profile", h.profileHandler, middleware.Auth(h.authConfig, h.authSvc))

}
