package productkeyboard

import (
	"fmt"
	"telegram-service-platform/delivery/telegramserver/callback"
	"telegram-service-platform/delivery/telegramserver/keyboard"
	"telegram-service-platform/params/productparams"

	"github.com/go-telegram/bot/models"
)

func PremiumPlans(response productparams.GetPremiumPlansResponse) *models.InlineKeyboardMarkup {

	buttons := make([]keyboard.Button, 0)
	for _, plan := range response.Plans {

		buttons = append(
			buttons,
			keyboard.Button{
				Text: fmt.Sprintf(
					"👑 Telegram Premium %d Month - %.2f USDT",
					plan.Months,
					plan.Price.USDT,
				),

				Data: fmt.Sprintf(
					"%s:%d",
					callback.ProductPremiumSelectCallBack,
					plan.ID,
				),

				Style: keyboard.Primary,
			})
	}

	builder := keyboard.NewBuilder()

	builder.AddButtonsPerRow(buttons, 1)

	builder.AddRow(keyboard.Back(callback.UserMainMenuCallBack))

	return builder.Build()
}
