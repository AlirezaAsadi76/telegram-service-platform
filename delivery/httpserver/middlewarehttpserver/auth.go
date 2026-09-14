package middlewarehttpserver

import (
	"net/http"
	"telegram-service-platform/config"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func (m Middleware) Auth() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey:    config.AuthMiddlewareContextKey,
		SigningKey:    m.authConfig.SignKey,
		SigningMethod: echojwt.AlgorithmHS256,
		ParseTokenFunc: func(c *echo.Context, auth string) (interface{}, error) {
			claims, cErr := m.authSvc.ParseToken(auth)
			if cErr != nil {
				return nil, cErr
			}
			return claims, nil
		},
		ErrorHandler: func(c *echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error":   "unauthorized",
				"message": "You are not authorized to access this resource",
			})
		},
	})
}
