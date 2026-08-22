package keyboard

import (
	"telegram-service-platform/delivery/telegramserver/callback"

	"github.com/go-telegram/bot/models"
)

// OrderConfirmMenu کیبورد تأیید نهایی سفارش را با استفاده از Builder می‌سازد
func OrderConfirmMenu() *models.InlineKeyboardMarkup {
	builder := NewBuilder()

	// دکمه پرداخت
	builder.AddRow(
		Button{Text: "💰 پرداخت با کیف پول", Data: callback.OrderPayWallet, Style: Success},
	)

	// دکمه انصراف
	builder.AddRow(
		Back(callback.OrderCancel),
	)

	return builder.Build()
}
