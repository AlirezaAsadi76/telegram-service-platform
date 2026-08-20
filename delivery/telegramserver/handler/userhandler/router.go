package userhandler

import (
	"github.com/go-telegram/bot"
)

func (h Handler) RegisterRoutes(b *bot.Bot) {

	//b.RegisterHandler(
	//	bot.HandlerTypeMessageText,
	//	"/start",
	//	bot.MatchTypePrefix,
	//	h.start,
	//	middleware.Public()...,
	//)
	//
	//b.RegisterHandler(
	//	bot.HandlerTypeCallbackQueryData,
	//	callback.UserMainMenuCallBack,
	//	bot.MatchTypePrefix,
	//	h.callback,
	//	middleware.Public()...,
	//)

}
