package producthandler

import (
	"context"
	"strings"
	"telegram-service-platform/delivery/telegramserver/callback"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) starsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.CallbackQuery == nil {
		return
	}
	data := update.CallbackQuery.Data

	switch {
	case strings.HasPrefix(data, callback.ProductStarsSelectCallBack):
		h.showStars(ctx, b, update) //TODO - must be complete
	case strings.HasPrefix(data, callback.ProductStarsBuyCallBack):
		h.showStars(ctx, b, update) //TODO - must be complete
	default:
		h.showStars(ctx, b, update)
	}

}
