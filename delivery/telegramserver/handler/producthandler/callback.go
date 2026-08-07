package producthandler

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) callback(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.CallbackQuery == nil {
		return
	}
	data := update.CallbackQuery.Data

	switch {
	case strings.HasPrefix(data, ProductStarsCallBack):
		h.showStars(ctx, b, update)
	case strings.HasPrefix(data, ProductPremiumCallBack):
		h.showPremium(ctx, b, update)
	}

}
