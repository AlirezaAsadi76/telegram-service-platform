package userhandler

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

func (h *Handler) SetRoutes(e *echo.Echo) {
	fmt.Println("SetRoutes UserHandler")
	userGroup := e.Group("/users")
	userGroup.POST("/register", h.registerHandler)
	userGroup.POST("/login", h.loginHandler)
	userGroup.GET("/balance", h.balanceHandler, h.middleware.Auth())
	userGroup.GET("/profile", h.profileHandler, h.middleware.Auth())
	userGroup.GET("/transactions", h.transactionHandler, h.middleware.Auth())

}
