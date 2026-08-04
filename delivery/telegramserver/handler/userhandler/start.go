package userhandler

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"telegram-service-platform/params"
)

func (h Handler) start(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {

	tgUser := update.Message.From

	request := params.GetOrRegisterRequest{

		TelegramID: tgUser.ID,
		Username:   tgUser.Username,
		FirstName:  tgUser.FirstName,
		LastName:   tgUser.LastName,
	}

	err := h.userValidator.GetOrRegister(request)
	if err != nil {
		log.Println(err)
		return
	}
	_, gErr := h.userService.GetOrRegister(
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
}
