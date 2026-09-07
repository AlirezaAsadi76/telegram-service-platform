package claims

import (
	"telegram-service-platform/config"
	"telegram-service-platform/service/authservice"

	"github.com/labstack/echo/v5"
)

func GetClaimsFromEchoContext(c *echo.Context) *authservice.Claims {
	rawClaims := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)
	return rawClaims
}
