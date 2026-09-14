package paymentflowtesting

import "github.com/labstack/echo/v5"

type MiddlewareTest struct{}

func newFakeMiddleware() MiddlewareTest {
	return MiddlewareTest{}
}

func (m MiddlewareTest) Auth() echo.MiddlewareFunc {
	return nil
}
