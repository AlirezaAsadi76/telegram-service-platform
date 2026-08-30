package keyboard

import (
	"telegram-service-platform/delivery/telegramserver/callback"

	"github.com/go-telegram/bot/models"
)

func OrderConfirmMenu() *models.InlineKeyboardMarkup {
	builder := NewBuilder()

	builder.AddRow(
		Button{
			Text:  "💰 پرداخت با کیف پول",
			Data:  callback.OrderPayWallet,
			Style: Success,
		},
	)

	builder.AddRow(
		Button{
			Text:  "💳 پرداخت با درگاه",
			Data:  callback.OrderPayGateway,
			Style: Success,
		},
	)

	builder.AddRow(
		Button{
			Text:  "🪙 پرداخت با رمز ارز",
			Data:  callback.OrderPayCrypto,
			Style: Success,
		},
	)

	builder.AddRow(
		Button{
			Text:  "❌ انصراف از سفارش",
			Data:  callback.OrderCancel,
			Style: Danger,
		},
	)

	return builder.Build()
}
