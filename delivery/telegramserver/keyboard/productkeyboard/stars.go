package productkeyboard

import (
	"fmt"

	"telegram-service-platform/delivery/telegramserver/callback"
	"telegram-service-platform/delivery/telegramserver/keyboard"
	"telegram-service-platform/params"

	"github.com/go-telegram/bot/models"
)

func StarsPlans(
	response params.GetStarPlansResponse,
) *models.InlineKeyboardMarkup {

	buttons := make([]keyboard.Button, 0)

	for _, plan := range response.Plans {

		buttons = append(
			buttons,
			keyboard.Button{
				Text: fmt.Sprintf(
					"%d ⭐ - %.2f USDT",
					plan.Amount,
					plan.Price.USDT,
				),
				Data: fmt.Sprintf(
					"%s:%d",
					callback.ProductStarsSelectCallBack,
					plan.ID,
				),
				Style: keyboard.Primary,
			},
		)
	}

	builder := keyboard.NewBuilder()

	builder.AddButtonsPerRow(
		buttons,
		2,
	)

	builder.AddRow(
		keyboard.Back(callback.UserMainMenuCallBack),
	)

	return builder.Build()
}
