package telegramserver

import "context"

func (b *Bot) Start(ctx context.Context) {

	b.client.Start(ctx)

}
