package keyboard

import (
	"telegram-service-platform/delivery/telegramserver/callback"

	"github.com/go-telegram/bot/models"
)

func PaymentActionMenu(paymentURL string) *models.InlineKeyboardMarkup {
	builder := NewBuilder()

	builder.AddRow(
		Button{
			Text:  "🌐 ورود به صفحه پرداخت",
			URL:   paymentURL,
			Style: Success,
		},
	)
	
	builder.AddRow(
		Button{
			Text:  "✅ بررسی وضعیت پرداخت",
			Data:  callback.OrderCheckPay,
			Style: Primary,
		},
	)

	builder.AddRow(
		Button{
			Text:  "❌ انصراف",
			Data:  callback.OrderCancel,
			Style: Danger,
		},
	)

	return builder.Build()
}
