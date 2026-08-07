package telegramserver

import "context"

func (b *Bot) Shutdown(ctx context.Context) error {

	// telegram library shutdown
	_, cErr := b.client.Close(ctx)
	return cErr
}
