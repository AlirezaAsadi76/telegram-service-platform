package producthandler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h Handler) buyStars(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {

	// TODO:
	// create order
	// call order service
	// send result

}
