package keyboard

import (
	"telegram-service-platform/delivery/telegramserver/callback"

	"github.com/go-telegram/bot/models"
)

func MainMenu() *models.InlineKeyboardMarkup {

	return &models.InlineKeyboardMarkup{

		InlineKeyboard: [][]models.InlineKeyboardButton{

			{
				{
					Text:              "Telegram Stars",
					CallbackData:      callback.ProductStarsCallBack,
					Style:             "primary",
					IconCustomEmojiID: "5235630047959727475",
				},
			},

			{
				{
					Text:         "👑 Telegram Premium",
					CallbackData: callback.ProductPremiumCallBack,
					Style:        "primary",
				},
			},

			{
				{
					Text:         "📢 Telegram Ads",
					CallbackData: callback.ProductAdsCallBack,
					Style:        "primary",
				},
			},

			{
				{
					Text:         "💎 Wallet",
					CallbackData: callback.Wallet,
					Style:        "success",
				},
			},

			{
				{
					Text:         "🎁 Gift",
					CallbackData: callback.Gift,
					Style:        "success",
				},
			},
		},
	}
}
