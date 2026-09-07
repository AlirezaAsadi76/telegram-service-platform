package userhandler

import (
	"net/http"
	"telegram-service-platform/entity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/userparams"
	"telegram-service-platform/pkg/telegramutils"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func (h *Handler) loginHandler(c *echo.Context) error {
	var req userparams.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	parsedData, vErr := telegramutils.VerifyTelegramInitData(req.InitData, h.botToken)
	if vErr != nil {
		logger.Logger.Warn("telegram login failed: invalid signature", zap.Error(vErr))
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid telegram data"})
	}

	telegramUser, pErr := telegramutils.ParseTelegramUser(parsedData.Get("user"))
	if pErr != nil {
		logger.Logger.Error("failed to parse telegram user json", zap.Error(pErr))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	userInfo, uErr := h.userSvc.GetOrRegister(c.Request().Context(), userparams.GetOrRegisterRequest{
		TelegramID: telegramUser.ID,
		Username:   telegramUser.Username,
		FirstName:  telegramUser.FirstName,
	})
	if uErr != nil {
		logger.Logger.Error("failed to get or register user", zap.Error(uErr))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to process user"})
	}

	user := entity.User{
		TelegramID: telegramUser.ID,
		Username:   telegramUser.Username,
		FirstName:  telegramUser.FirstName,
		LastName:   telegramUser.LastName,
		Role:       userInfo.UserInfo.Role,
	}

	accessToken, aErr := h.authSvc.CreateAccessToken(user)
	if aErr != nil {
		logger.Logger.Error("failed to create access token", zap.Error(aErr))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	refreshToken, rErr := h.authSvc.CreateRefreshToken(user)
	if rErr != nil {
		logger.Logger.Error("failed to create refresh token", zap.Error(rErr))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, userparams.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}
