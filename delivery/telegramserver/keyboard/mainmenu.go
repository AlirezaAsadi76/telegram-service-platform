package keyboard

import (
	"telegram-service-platform/delivery/telegramserver/handler/callbackhandler"

	"github.com/go-telegram/bot/models"
)

func MainMenu() *models.InlineKeyboardMarkup {

	return &models.InlineKeyboardMarkup{

		InlineKeyboard: [][]models.InlineKeyboardButton{

			{
				{
					Text:         "Telegram Stars",
					CallbackData: callbackhandler.Stars,
					Style:        "primary",
				},
			},

			{
				{
					Text:         "👑 Telegram Premium",
					CallbackData: callbackhandler.Premium,
					Style:        "primary",
				},
			},

			{
				{
					Text:         "📢 Telegram Ads",
					CallbackData: callbackhandler.Ads,
					Style:        "primary",
				},
			},

			{
				{
					Text:         "💎 Wallet",
					CallbackData: callbackhandler.Wallet,
					Style:        "success",
				},
			},

			{
				{
					Text:         "🎁 Gift",
					CallbackData: callbackhandler.Gift,
					Style:        "success",
				},
			},
		},
	}
}
