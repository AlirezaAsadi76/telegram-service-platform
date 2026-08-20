package keyboard

import (
	"fmt"
	"telegram-service-platform/delivery/telegramserver/callback"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/helpers"

	"github.com/go-telegram/bot/models"
)

func MainMenu(platforms []smmentity.Platform, categories []smmentity.Category) *models.InlineKeyboardMarkup {
	builder := NewBuilder()
	fmt.Println(platforms)
	fmt.Println(categories)
	builder.AddRow(Button{
		Text:  "📢 تبلیغات هوشمند تلگرام",
		Data:  callback.MainMenuTelegramAds,
		Style: Primary,
	})

	builder.AddRow(
		Button{Text: "⭐ خرید استارز", Data: callback.MainMenuStars, Style: Primary},
		Button{Text: "💎 خرید پریمیوم", Data: callback.MainMenuPremium, Style: Primary},
	)

	var catButtons []Button
	for _, cat := range categories {
		catButtons = append(catButtons, Button{
			Text:  fmt.Sprintf("%s %s", helpers.GetCategoryIcon(cat.Name.String()), helpers.GetCategoryDisplayName(cat.Name.String())),
			Data:  fmt.Sprintf("%s:%s:%s", callback.SMMCategoryPrefix, "telegram", cat.Name),
			Style: Primary,
		})
	}
	if len(catButtons) > 0 {
		builder.AddButtonsPerRow(catButtons, 2)
	}

	// ردیف ۵: سایر پلتفرم‌ها (از دیتابیس، حداکثر ۲ تا در هر ردیف برای خوانایی بهتر در موبایل)
	var platformButtons []Button
	for _, platform := range platforms {
		if platform.Name == "telegram" {
			continue
		}
		platformButtons = append(platformButtons, Button{
			Text:  fmt.Sprintf("%s %s", helpers.GetPlatformIcon(platform.Name.String()), helpers.GetPlatformDisplayName(platform.Name.String())),
			Data:  fmt.Sprintf("%s:%s", callback.SMMPlatformPrefix, platform.Name),
			Style: Primary,
		})
	}
	if len(platformButtons) > 0 {
		builder.AddButtonsPerRow(platformButtons, 4)
	}

	// ردیف ۶: کیف پول و قوانین
	builder.AddRow(
		Button{Text: "💰 کیف پول", Data: callback.MainMenuWallet, Style: Success},
		Button{Text: "📜 قوانین", Data: callback.MainMenuRules, Style: Primary},
	)

	return builder.Build()
}
