package producthandler

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) selectStars(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {

	data := update.CallbackQuery.Data

	parts := strings.Split(data, ":")

	if len(parts) != 4 {
		return
	}

	id, err := strconv.ParseInt(parts[3], 10, 64)

	if err != nil {
		return
	}

	_ = id

	// TODO:
	// get plan detail
	// create confirm keyboard
	// edit message

}
