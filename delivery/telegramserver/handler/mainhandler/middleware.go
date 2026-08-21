package mainhandler

import (
	"telegram-service-platform/delivery/telegramserver/middleware"

	"github.com/go-telegram/bot"
)

func (h *Handler) middlewares() []bot.Middleware {
	return append(
		middleware.Public(),
		middleware.ActivityTracking(h.userService),
	)
}
