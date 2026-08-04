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
					Text:         "Telegram Stars",
					CallbackData: callback.Stars,
					Style:        "primary",
				},
			},

			{
				{
					Text:         "👑 Telegram Premium",
					CallbackData: callback.Premium,
					Style:        "primary",
				},
			},

			{
				{
					Text:         "📢 Telegram Ads",
					CallbackData: callback.Ads,
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
