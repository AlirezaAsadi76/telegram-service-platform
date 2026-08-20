package keyboard

import (
	"github.com/go-telegram/bot/models"
)

// ReplyMainMenu builds the permanent reply keyboard (always visible at bottom)
func ReplyMainMenu() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{
					Text: "👤 حساب کاربری",
				},
				{
					Text: "💳 تراکنش‌ها",
				},
			},
			{
				{
					Text: "🏠 منوی اصلی",
				},
				{
					Text: " ارتباط با ادمین",
				},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

// ReplyMainMenuTexts returns the text labels for reply keyboard buttons
// Used for matching in router.go
func ReplyMainMenuTexts() []string {
	return []string{
		"👤 حساب کاربری",
		"💳 تراکنش‌ها",
		"🏠 منوی اصلی",
		"📞 ارتباط با ادمین",
	}
}
