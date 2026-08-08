package producthandler

import (
	"context"
	"strings"
	"telegram-service-platform/delivery/telegramserver/callback"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) premiumHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.CallbackQuery == nil {
		return
	}
	data := update.CallbackQuery.Data

	switch {
	case strings.HasPrefix(data, callback.ProductPremiumSelectCallBack):
		h.showPremium(ctx, b, update) //TODO - must be complete
	case strings.HasPrefix(data, callback.ProductPremiumBuyCallBack):
		h.showPremium(ctx, b, update) //TODO - must be complete
	default:
		h.showPremium(ctx, b, update)

	}

}
