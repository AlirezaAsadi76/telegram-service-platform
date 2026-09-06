package middleware

import (
	"net/http"
	"telegram-service-platform/config"
	"telegram-service-platform/service/authservice"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func Auth(authConfig authservice.Config, authService *authservice.Service) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey:    config.AuthMiddlewareContextKey,
		SigningKey:    authConfig.SignKey,
		SigningMethod: echojwt.AlgorithmHS256,
		ParseTokenFunc: func(c *echo.Context, auth string) (interface{}, error) {
			claims, cErr := authService.ParseToken(auth)
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
