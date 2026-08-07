package producthandler

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) showStars(ctx context.Context, b *bot.Bot, update *models.Update) {

	response, sErr := h.productService.GetStarPlans(ctx)
	if sErr != nil {
		log.Println(sErr)
	}

	log.Println(response)
}
