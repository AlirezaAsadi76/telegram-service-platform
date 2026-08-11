package userhandler

import (
	"context"
	"log"
	"telegram-service-platform/delivery/telegramserver/keyboard"
	"telegram-service-platform/pkg/mapper"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) start(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.Message == nil {
		return
	}
	request := mapper.MapTelegramUserToRegisterRequest(update.Message.From)

	err := h.userValidator.GetOrRegister(request)
	//TODO - create users wallet
	if err != nil {
		log.Println(err)
		return
	}

	user, gErr := h.userService.GetOrRegister(ctx, request)
	if gErr != nil {
		return
	}

	merr := h.messenger.Send(
		ctx,
		b,
		&bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "به پنل خدمات تلگرام خوش آمدید 👋",
			ReplyMarkup: keyboard.MainMenu(),
		},
	)

	if merr != nil {
		log.Println(merr)
		return
	}

	log.Println(
		"registered user:",
		user.UserInfo,
	)
}
