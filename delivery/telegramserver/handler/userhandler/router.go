package userhandler

import (
	"telegram-service-platform/delivery/telegramserver/callback"
	"telegram-service-platform/delivery/telegramserver/middleware"

	"github.com/go-telegram/bot"
)

func (h Handler) RegisterRoutes(b *bot.Bot) {

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.MainMenuWallet,
		bot.MatchTypeExact,
		h.WalletBalance,
		middleware.Public()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		callback.ReplyTransactionsText,
		bot.MatchTypeExact,
		h.TransactionsHistory,
		middleware.Public()...,
	)

}
