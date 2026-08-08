package telegramserver

import "log"

func (b *Bot) registerRoutes() {
	log.Println("register routes:", len(b.handlers))
	for _, h := range b.handlers {
		h.RegisterRoutes(b.client)
	}
}
