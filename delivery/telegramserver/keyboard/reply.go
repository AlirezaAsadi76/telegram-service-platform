package keyboard

import "github.com/go-telegram/bot/models"

func ReplyMainMenu() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "👤 حساب کاربری"},
				{Text: "💳 تراکنش‌ها"},
			},
			{
				{Text: "🏠 منوی اصلی"},
				{Text: "📞 ارتباط با ادمین"},
			},
		},
		ResizeKeyboard: true,
		IsPersistent:   true,
	}
}
