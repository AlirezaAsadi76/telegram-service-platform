package userhandler

import (
	"fmt"
	"telegram-service-platform/delivery/httpserver/middleware"

	"github.com/labstack/echo/v5"
)

func (h *Handler) SetRoutes(e *echo.Echo) {
	fmt.Println("SetRoutes UserHandler")
	userGroup := e.Group("/users")
	userGroup.POST("/register", h.registerHandler)
	userGroup.POST("/login", h.loginHandler)
	userGroup.GET("/balance", h.balanceHandler, middleware.Auth(h.authConfig, h.authSvc))
	userGroup.GET("/profile", h.profileHandler, middleware.Auth(h.authConfig, h.authSvc))
	userGroup.GET("/transactions", h.transactionHandler, middleware.Auth(h.authConfig, h.authSvc))

}
