package producthandler

import (
	"context"
	"log"

	"telegram-service-platform/delivery/telegramserver/keyboard/productkeyboard"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) showStars(ctx context.Context, b *bot.Bot, update *models.Update) {

	response, err := h.productService.GetStarPlans(ctx)

	if err != nil {
		log.Println("get stars plans:", err)
		return
	}

	err = h.messenger.Send(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "لطفا تعداد Telegram Stars را انتخاب کنید:",
		ReplyMarkup: productkeyboard.StarsPlans(
			response,
		),
	},
	)

	if err != nil {
		log.Println("send stars plans:", err)
	}

}
