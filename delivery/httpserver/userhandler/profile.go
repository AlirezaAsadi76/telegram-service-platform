package userhandler

import (
	"fmt"
	"net/http"
	"telegram-service-platform/entity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/userparams"
	claimspkg "telegram-service-platform/pkg/claims"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func (h *Handler) profileHandler(c *echo.Context) error {
	const op = "userhandler.profileHandler"
	fmt.Println(op)
	claims := claimspkg.GetClaimsFromEchoContext(c)
	fmt.Println("claims: ", claims)
	user, fErr := h.userSvc.FindUserByTelegramID(c.Request().Context(), userparams.FindUserByTelegramIDRequest{
		TelegramID: entity.TelegramId(claims.TelegramId),
	})
	if fErr != nil {
		logger.Logger.Error("failed to get user by id",
			zap.String("op", op),
			zap.Int64("user_id", claims.UserId),
			zap.Error(fErr),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch user profile"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":     user.UserInfo.Id,
		"telegram_id": user.UserInfo.TelegramID,
		"username":    user.UserInfo.Username,
		"first_name":  user.UserInfo.FirstName,
		"role":        user.UserInfo.Role,
		"created_at":  user.UserInfo.CreatedAt,
	})

}
