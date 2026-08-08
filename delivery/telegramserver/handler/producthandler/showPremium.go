package producthandler

import (
	"context"
	"log"

	"telegram-service-platform/delivery/telegramserver/keyboard/productkeyboard"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) showPremium(ctx context.Context, b *bot.Bot, update *models.Update) {

	response, err := h.productService.GetPremiumPlans(ctx)

	if err != nil {
		log.Println("get premium plans:", err)
		return
	}

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,

			Text: "مدت زمان Telegram Premium را انتخاب کنید:",

			ReplyMarkup: productkeyboard.PremiumPlans(
				response,
			),
		},
	)

	if err != nil {
		log.Println("send premium plans:", err)
	}

}
