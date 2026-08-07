package handler

import "github.com/go-telegram/bot"

type Handler interface {
	RegisterRoutes(b *bot.Bot)
}
