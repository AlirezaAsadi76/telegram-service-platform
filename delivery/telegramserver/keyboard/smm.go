package keyboard

import (
	"fmt"
	"telegram-service-platform/delivery/telegramserver/callback"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/helpers"

	"github.com/go-telegram/bot/models"
)

func CategoryMenu(platform string, categories []smmentity.Category) *models.InlineKeyboardMarkup {
	builder := NewBuilder()

	var buttons []Button
	for _, cat := range categories {
		buttons = append(buttons, Button{
			Text:  fmt.Sprintf("%s %s", helpers.GetCategoryIcon(cat.Name.String()), helpers.GetCategoryDisplayName(cat.Name.String())),
			Data:  fmt.Sprintf("%s:%s:%s", callback.SMMCategoryPrefix, platform, cat.Name),
			Style: Primary,
		})
	}

	if len(buttons) > 0 {
		builder.AddButtonsPerRow(buttons, 2)
	}

	builder.AddRow(Back(callback.UserMainMenuCallBack))
	return builder.Build()
}

func ServiceMenu(platform, category string, services []smmentity.SmmMapping) *models.InlineKeyboardMarkup {
	builder := NewBuilder()

	var buttons []Button
	for _, svc := range services {
		buttons = append(buttons, Button{
			Text:  svc.ButtonName,
			Data:  fmt.Sprintf("%s:%s:%s:%d", callback.SMMServicePrefix, platform, category, svc.Id),
			Style: Primary,
		})
	}

	if len(buttons) > 0 {

		builder.AddButtonsPerRow(buttons, 1)
	}

	builder.AddRow(Back(fmt.Sprintf("%s:%s", callback.SMMPlatformPrefix, platform)))
	return builder.Build()
}
