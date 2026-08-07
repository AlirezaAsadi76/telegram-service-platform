package telegramserver

import (
	"context"
	"log"
)

func (b *Bot) Start(ctx context.Context) error {

	log.Println("telegram server started")

	b.client.Start(ctx)

	log.Println("telegram server stopped")

	return nil
}
