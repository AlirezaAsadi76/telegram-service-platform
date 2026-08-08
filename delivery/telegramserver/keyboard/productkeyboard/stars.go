package productkeyboard

import (
	"fmt"
	"telegram-service-platform/delivery/telegramserver/callback"

	"telegram-service-platform/params"

	"github.com/go-telegram/bot/models"
)

func StarsPlans(response params.GetStarPlansResponse) *models.InlineKeyboardMarkup {

	buttons := make([][]models.InlineKeyboardButton, 0, len(response.Plans))

	for _, plan := range response.Plans {
		buttons = append(
			buttons,
			[]models.InlineKeyboardButton{
				{
					Text: fmt.Sprintf(
						"%d ⭐ - %.2f USDT",
						plan.Amount,
						plan.Price.USDT,
					),
					CallbackData: fmt.Sprintf(
						"%s:%d",
						callback.ProductStarsSelectCallBack,
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
