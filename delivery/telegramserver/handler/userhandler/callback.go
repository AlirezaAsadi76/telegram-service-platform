package userhandler

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) callback(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.CallbackQuery == nil {
		return
	}
	log.Println(
		"callback:", update.CallbackQuery.Data)

}
