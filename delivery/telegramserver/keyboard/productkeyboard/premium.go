package productkeyboard

import (
	"fmt"
	"telegram-service-platform/delivery/telegramserver/callback"

	"telegram-service-platform/params"

	"github.com/go-telegram/bot/models"
)

func PremiumPlans(response params.GetPremiumPlansResponse) *models.InlineKeyboardMarkup {

	buttons := make([][]models.InlineKeyboardButton, 0, len(response.Plans))

	for _, plan := range response.Plans {

		buttons = append(
			buttons,
			[]models.InlineKeyboardButton{
				{
					Text: fmt.Sprintf(
						"👑 Telegram Premium %d Month - %.2f USDT",
						plan.Months,
						plan.Price.USDT,
					),

					CallbackData: fmt.Sprintf(
						"%s:%d",
						callback.ProductPremiumSelectCallBack,
						plan.ID,
					),

					Style: "primary",
				},
			},
		)
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}
