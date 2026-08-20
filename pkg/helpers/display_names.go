package helpers

import "strings"

func GetPlatformDisplayName(platform string) string {
	platformLoverCase := strings.ToLower(platform)

	platforms := map[string]string{
		"telegram":  "تلگرام",
		"instagram": "اینستاگرام",
		"tiktok":    "تیک‌تاک",
		"whatsapp":  "واتساپ",
		"twitter":   "توییتر",
	}

	if name, ok := platforms[platformLoverCase]; ok {
		return name
	}
	return platform
}

func GetCategoryDisplayName(category string) string {
	categories := map[string]string{
		"views":     "بازدید",
		"reactions": "ری‌اکشن",
		"members":   "ممبر",
		"shares":    "فوروارد / شیر",
		"likes":     "لایک",
		"followers": "فالوور",
		"comments":  "کامنت",
	}

	categoryLoverCase := strings.ToLower(category)

	if name, ok := categories[categoryLoverCase]; ok {
		return name
	}
	return category
}

func GetCategoryIcon(category string) string {
	categoryLoverCase := strings.ToLower(category)
	icons := map[string]string{
		"views":     "👁️",
		"reactions": "❤️",
		"members":   "👥",
		"shares":    "🔄",
		"likes":     "👍",
		"followers": "👤",
		"comments":  "💬",
	}

	if icon, ok := icons[categoryLoverCase]; ok {
		return icon
	}
	return "📦"
}

// GetPlatformIcon returns emoji for platform
func GetPlatformIcon(platform string) string {
	platformLoverCase := strings.ToLower(platform)
	icons := map[string]string{
		"telegram":  "✈️",
		"instagram": "📷",
		"tiktok":    "🎵",
		"whatsapp":  "💬",
		"twitter":   "𝕏",
	}

	if icon, ok := icons[platformLoverCase]; ok {
		return icon
	}
	return "🌐"
}
