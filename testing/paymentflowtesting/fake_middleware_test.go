package paymentflowtesting

import "github.com/labstack/echo/v5"

type MiddlewareTest struct{}

func newFakeMiddleware() MiddlewareTest {
	return MiddlewareTest{}
}

func (m MiddlewareTest) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return next
	}
}
