package userhandler

import (
	"context"
	"log"
	"telegram-service-platform/pkg/mapper"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) start(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	log.Println("userHandler start")
	if update.Message == nil {
		return
	}
	request := mapper.MapTelegramUserToRegisterRequest(update.Message.From)

	err := h.userValidator.GetOrRegister(request)
	if err != nil {
		log.Println(err)
		return
	}
	user, gErr := h.userService.GetOrRegister(
		ctx,
		request,
	)
	if gErr != nil {
		log.Println(gErr)
		return
	}

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "خوش آمدید 👋",
		},
	)

	if err != nil {
		log.Println(err)
		return
	}
	log.Println(
		"registered user:",
		user.UserInfo,
	)
}
