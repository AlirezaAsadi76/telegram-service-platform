package middleware

import (
	"context"
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

			log.Printf(
				"telegram update received: %d",
				update.ID,
			)

			next(ctx, b, update)

			log.Printf(
				"telegram update completed in %s",
				time.Since(start),
			)

		}

	}
}
