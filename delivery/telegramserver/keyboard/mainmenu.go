package keyboard

import (
	"github.com/go-telegram/bot/models"
)

func MainMenu() *models.InlineKeyboardMarkup {

	return &models.InlineKeyboardMarkup{

		InlineKeyboard: [][]models.InlineKeyboardButton{

			{
				{
					Text:         "⭐ Telegram Stars",
					CallbackData: "stars",
				},
			},

			{
				{
					Text:         "👑 Telegram Premium",
					CallbackData: "premium",
				},
			},

			{
				{
					Text:         "📢 Telegram Ads",
					CallbackData: "ads",
				},
			},

			{
				{
					Text:         "💎 Wallet",
					CallbackData: "wallet",
				},
			},

			{
				{
					Text:         "🎁 Gift",
					CallbackData: "gift",
				},
			},
		},
	}
}
