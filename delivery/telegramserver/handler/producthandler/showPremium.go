package producthandler

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) showPremium(ctx context.Context, b *bot.Bot, update *models.Update) {

	response, pErr := h.productService.GetPremiumPlans(ctx)
	if pErr != nil {
		log.Println(pErr)
	}
	log.Println(response)
}
