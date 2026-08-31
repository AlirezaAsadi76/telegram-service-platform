package orderhandler

import (
	"github.com/go-telegram/bot"

	"telegram-service-platform/delivery/telegramserver/callback"
	"telegram-service-platform/delivery/telegramserver/middleware"
)

func (h *Handler) middlewares() []bot.Middleware {
	return append(middleware.Public(), middleware.ActivityTracking(h.userService)) // یا هر سرویس کاربری که دارید
}

func (h *Handler) RegisterRoutes(b *bot.Bot) {

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.OrderPayWallet,
		bot.MatchTypeExact,
		h.processWalletPayment,
		h.middlewares()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.OrderCancel,
		bot.MatchTypeExact,
		h.cancelOrder,
		h.middlewares()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.OrderCheckPay,
		bot.MatchTypeExact,
		h.checkPaymentStatus,
		middleware.Public()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.OrderPayCrypto,
		bot.MatchTypeExact,
		h.processCryptoPayment,
		middleware.Public()...,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		callback.OrderPayGateway,
		bot.MatchTypeExact,
		h.processGatewayPayment,
		middleware.Public()...,
	)

}
