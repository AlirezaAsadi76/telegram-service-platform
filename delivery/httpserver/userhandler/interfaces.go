package userhandler

import "github.com/labstack/echo/v5"

type AuthMiddleware interface {
	Auth() echo.MiddlewareFunc
}
