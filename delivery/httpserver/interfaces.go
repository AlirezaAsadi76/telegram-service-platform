package httpserver

import "github.com/labstack/echo/v5"

type Handlers interface {
	SetRoutes(e *echo.Echo)
}
