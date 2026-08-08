package userhandler

import (
	"context"
	"strings"
	"telegram-service-platform/delivery/telegramserver/callback"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) callback(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.CallbackQuery == nil {
		return
	}
	data := update.CallbackQuery.Data

	switch {
	case strings.HasPrefix(data, callback.UserMainMenuCallBack):
		h.start(ctx, b, update)

	}

}
