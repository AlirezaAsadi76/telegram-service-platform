package userhandler

import (
	"net/http"
	"telegram-service-platform/params/userparams"

	"github.com/labstack/echo/v5"
)

func (h *Handler) registerHandler(c *echo.Context) error {
	const op = "userhandler.register"
	var registerRequest userparams.GetOrRegisterRequest

	if err := c.Bind(&registerRequest); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{})
	}

	if err := h.userVal.GetOrRegister(registerRequest); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{})
	}

	res, gErr := h.userSvc.GetOrRegister(c.Request().Context(), registerRequest)
	if gErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{})
	}

	return c.JSON(http.StatusOK, res)
}
