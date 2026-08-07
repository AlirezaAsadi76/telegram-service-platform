package keyboard

import (
	"telegram-service-platform/delivery/telegramserver/handler/producthandler"

	"github.com/go-telegram/bot/models"
)

func MainMenu() *models.InlineKeyboardMarkup {

	return &models.InlineKeyboardMarkup{

		InlineKeyboard: [][]models.InlineKeyboardButton{

			{
				{
					Text:         "Telegram Stars",
					CallbackData: producthandler.ProductStarsCallBack,
					Style:        "primary",
				},
			},

			{
				{
					Text:         "👑 Telegram Premium",
					CallbackData: producthandler.ProductPremiumCallBack,
					Style:        "primary",
				},
			},

			{
				{
					Text:         "📢 Telegram Ads",
					CallbackData: producthandler.ProductAdsCallBack,
					Style:        "primary",
				},
			},

			{
				{
					Text:         "💎 Wallet",
					CallbackData: producthandler.Wallet,
					Style:        "success",
				},
			},

			{
				{
					Text:         "🎁 Gift",
					CallbackData: producthandler.Gift,
					Style:        "success",
				},
			},
		},
	}
}
