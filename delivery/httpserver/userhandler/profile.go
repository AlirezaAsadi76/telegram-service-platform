package userhandler

import (
	claimspkg "telegram-service-platform/pkg/claims"

	"github.com/labstack/echo/v5"
)

func (h *Handler) profileHandler(c *echo.Context) error {
	const op = "userhandler.profileHandler"
	claims := claimspkg.GetClaimsFromEchoContext(c)

	h.userSvc.
	return nil
}
