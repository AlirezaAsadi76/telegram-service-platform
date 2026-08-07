package producthandler

import (
	"log"

	"github.com/go-telegram/bot"
)

func (h Handler) RegisterRoutes(b *bot.Bot) {
	log.Println("product handler")

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		ProductPrefixCallBack,
		bot.MatchTypePrefix,
		h.callback,
	)

}
