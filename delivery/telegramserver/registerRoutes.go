package telegramserver

func (b *Bot) registerRoutes() {

	for _, h := range b.handlers {
		h.RegisterRoutes(b.client)
	}
}
