package producthandler

import (
	"telegram-service-platform/delivery/telegramserver/callback"
	"telegram-service-platform/delivery/telegramserver/middleware"

	"github.com/go-telegram/bot"
)

func (h Handler) RegisterRoutes(b *bot.Bot) {

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.ProductPrefixCallBack,
		bot.MatchTypePrefix,
		h.callback,
		middleware.Public()...,
	)

}
