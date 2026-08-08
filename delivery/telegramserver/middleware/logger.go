package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Logger() Middleware {

	return func(next HandlerFunc) HandlerFunc {

		return func(
			ctx context.Context,
			b *bot.Bot,
			update *models.Update,
		) {

			start := time.Now()

			var data string
			if update.CallbackQuery != nil {

				data = fmt.Sprintf("callback query: %s", update.CallbackQuery.Data)
			}
			if update.Message != nil {
				data = fmt.Sprintf("message: %s", update.Message.Text)
			}
			log.Printf(
				"telegram update received: %d - %s",
				update.ID,
				data,
			)

			log.Printf(
				"telegram update completed in %s",
				time.Since(start),
			)
			next(ctx, b, update)
		}

	}
}
