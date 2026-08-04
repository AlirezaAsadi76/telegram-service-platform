package telegramserver

func (b *Bot) registerRoutes() {

	b.userHandler.RegisterRoutes(
		b.client,
	)
}
