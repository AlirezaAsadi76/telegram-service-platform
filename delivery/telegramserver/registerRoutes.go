package telegramserver

func (b *Bot) registerRoutes() {

	b.userRouter.Register(
		b.client,
	)

}
