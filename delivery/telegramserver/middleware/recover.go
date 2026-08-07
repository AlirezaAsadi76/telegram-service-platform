package middleware

import (
	"context"
	"log"
	"runtime/debug"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Recover() Middleware {

	return func(next HandlerFunc) HandlerFunc {

		return func(
			ctx context.Context,
			b *bot.Bot,
			update *models.Update,
		) {

			defer func() {

				if r := recover(); r != nil {

					log.Printf(
						"panic recovered: %v\n%s",
						r,
						debug.Stack(),
					)

				}

			}()

			next(ctx, b, update)

		}

	}
}
