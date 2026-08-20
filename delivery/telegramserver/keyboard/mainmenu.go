package keyboard

import (
	"github.com/go-telegram/bot/models"
	"telegram-service-platform/delivery/telegramserver/callback"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/helpers"
)

// MainMenu builds the main inline keyboard with 6 rows
func MainMenu(platforms []smmentity.PlatformType, categories []smmentity.CategoryType) *models.InlineKeyboardMarkup {
	var rows [][]models.InlineKeyboardButton

	// Row 1: Telegram Ads (full width)
	rows = append(rows, []models.InlineKeyboardButton{
		{
			Text:         "📢 تبلیغات هوشمند تلگرام",
			CallbackData: callback.MainMenuTelegramAds,
		},
	})

	// Row 2: Stars & Premium
	rows = append(rows, []models.InlineKeyboardButton{
		{
			Text:         "⭐ خرید استارز",
			CallbackData: callback.MainMenuStars,
		},
		{
			Text:         "💎 خرید پریمیوم",
			CallbackData: callback.MainMenuPremium,
		},
	})

	// Row 3 & 4: Telegram Categories (Views, Reactions, Members, Shares)
	// These come from database (smm_service_mappings WHERE platform='telegram')
	if len(categories) > 0 {
		var row3, row4 []models.InlineKeyboardButton
		for i, cat := range categories {
			btn := models.InlineKeyboardButton{
				Text:         helpers.GetCategoryIcon(cat.String()) + " " + helpers.GetCategoryDisplayName(cat.String()),
				CallbackData: callback.SMMCategoryPrefix + ":" + cat.String(),
			}
			if i < 2 {
				row3 = append(row3, btn)
			} else {
				row4 = append(row4, btn)
			}
		}
		if len(row3) > 0 {
			rows = append(rows, row3)
		}
		if len(row4) > 0 {
			rows = append(rows, row4)
		}
	}

	// Row 5: Other Platforms (Instagram, TikTok, WhatsApp, Twitter)
	if len(platforms) > 0 {
		var row []models.InlineKeyboardButton
		for _, platform := range platforms {
			if platform == "telegram" {
				continue // Skip telegram, already shown
			}
			btn := models.InlineKeyboardButton{
				Text:         helpers.GetPlatformIcon(platform.String()) + " " + helpers.GetPlatformDisplayName(platform.String()),
				CallbackData: callback.SMMPlatformPrefix + ":" + platform.String(),
			}
			row = append(row, btn)
		}
		if len(row) > 0 {
			rows = append(rows, row)
		}
	}

	// Row 6: Wallet & Rules
	rows = append(rows, []models.InlineKeyboardButton{
		{
			Text:         "💰 کیف پول",
			CallbackData: callback.MainMenuWallet,
		},
		{
			Text:         "📜 قوانین",
			CallbackData: callback.MainMenuRules,
		},
	})

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}
