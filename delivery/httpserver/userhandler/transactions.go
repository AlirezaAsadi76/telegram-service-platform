package userhandler

import (
	"net/http"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/orderparams"
	claimspkg "telegram-service-platform/pkg/claims"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func (h *Handler) transactionHandler(c *echo.Context) error {
	const op = "userhandler.transactionHandler"
	claim := claimspkg.GetClaimsFromEchoContext(c)

	ordersHistory, oErr := h.orderSvc.GetByUserID(c.Request().Context(), orderparams.GetByUserIdRequest{
		UserID: uint64(claim.UserId),
	})

	if oErr != nil {

		logger.Logger.Error("failed to get order history", zap.Error(oErr),
			zap.String("op", op),
			zap.Int64("user_id", claim.UserId),
			zap.Error(oErr),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get order history"})

	}

	return c.JSON(http.StatusOK, ordersHistory)

}
