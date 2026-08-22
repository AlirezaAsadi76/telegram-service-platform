package keyboard

import (
	"github.com/go-telegram/bot/models"
	"telegram-service-platform/delivery/telegramserver/callback"
)

// OrderConfirmMenu کیبورد تأیید نهایی سفارش را با استفاده از Builder می‌سازد
func OrderConfirmMenu() *models.InlineKeyboardMarkup {
	builder := NewBuilder()

	// دکمه پرداخت
	builder.AddRow(
		Button{Text: "💰 پرداخت با کیف پول", Data: callback.or, Style: Success},
	)

	// دکمه انصراف
	builder.AddRow(
		Back(callback.OrderCancel),
	)

	return builder.Build()
}
