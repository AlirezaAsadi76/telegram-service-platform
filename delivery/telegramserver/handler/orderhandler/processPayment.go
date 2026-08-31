package orderhandler

import (
	"context"
	"telegram-service-platform/entity/paymententity"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) processGatewayPayment(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "orderhandler.processGatewayPayment"

	h.processDirectPayment(ctx, b, update, paymententity.PaymentMethodZarinpal, op)
}

func (h *Handler) processCryptoPayment(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "orderhandler.processCryptoPayment"
	h.processDirectPayment(ctx, b, update, paymententity.PaymentMethodCrypto, op)
}
