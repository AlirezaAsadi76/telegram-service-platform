package userhandler

import (
	"net/http"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/walletparam"
	claimspkg "telegram-service-platform/pkg/claims"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func (h *Handler) balanceHandler(c *echo.Context) error {
	const op = "userhandler.balanceHandler"
	claims := claimspkg.GetClaimsFromEchoContext(c)
	balanceRes, bErr := h.walletSvc.GetBalance(c.Request().Context(), walletparam.GetBalanceRequest{
		UserID: uint64(claims.UserId),
	})

	if bErr != nil {
		logger.Logger.Error("failed to get wallet balance",
			zap.String("op", op),
			zap.Int64("user_id", claims.UserId),
			zap.Error(bErr),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch balance"})
	}
	return c.JSON(http.StatusOK, balanceRes)
}
