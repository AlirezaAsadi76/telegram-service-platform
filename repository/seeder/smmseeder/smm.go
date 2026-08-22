package smmseeder

import (
	"context"
	"fmt"

	"telegram-service-platform/repository/postgres"
)

// SeedSMMData داده‌های اولیه SMM را در دیتابیس درج می‌کند
// این متد idempotent است (با ON CONFLICT DO NOTHING)
func SeedSMMData(ctx context.Context, db *postgres.DB) error {
	conn := db.Connection()

	// ۱. درج سرویس‌های خام JAP (شبیه‌سازی AllServices)
	rawServices := []struct {
		ServiceID    int64
		Name         string
		Type         string
		Rate         string
		MinQty       int64
		MaxQty       int64
		DripFeed     bool
		Refill       bool
		Cancel       bool
		Category     string
		ProviderName string
	}{
		// Telegram Services
		{10010000890098, "Telegram Post Views 1K", "Default", "0.50", 100, 100000, false, true, true, "telegram_views", "justanotherpanel"},
		{10020000890098, "Telegram Post Views 5K", "Default", "2.20", 500, 500000, false, true, true, "telegram_views", "justanotherpanel"},
		{10030000890098, "Telegram Story Views", "Default", "0.30", 50, 50000, false, false, true, "telegram_views", "justanotherpanel"},
		{10040000890098, "Telegram Positive Reaction 👍", "Default", "1.50", 10, 10000, false, false, false, "telegram_reactions", "justanotherpanel"},
		{10050000890098, "Telegram Real Members", "Default", "15.00", 100, 50000, false, true, true, "telegram_members", "justanotherpanel"},
		{10060000890098, "Telegram Premium Members", "Default", "45.00", 50, 10000, false, true, true, "telegram_members", "justanotherpanel"},
		{10070000890098, "Telegram Fake Members (Cheap)", "Default", "3.50", 100, 100000, false, false, false, "telegram_members", "justanotherpanel"},
		{10080000890098, "Telegram Post Shares/Forwards", "Default", "2.00", 50, 20000, false, false, true, "telegram_shares", "justanotherpanel"},

		// Instagram Services
		{2001, "Instagram Likes - Real", "Default", "0.80", 50, 100000, false, true, true, "instagram_likes", "justanotherpanel"},
		{2002, "Instagram Likes - High Quality", "Default", "1.20", 100, 50000, false, true, true, "instagram_likes", "justanotherpanel"},
		{2003, "Instagram Followers - Real", "Default", "8.50", 100, 50000, false, true, true, "instagram_followers", "justanotherpanel"},
		{2004, "Instagram Post Views", "Default", "0.10", 100, 1000000, false, false, true, "instagram_views", "justanotherpanel"},

		// TikTok Services
		{30010000890098, "TikTok Views", "Default", "0.05", 100, 5000000, false, false, true, "tiktok_views", "justanotherpanel"},
		{30020000890098, "TikTok Likes", "Default", "1.00", 50, 100000, false, true, true, "tiktok_likes", "justanotherpanel"},
		{30030000890098, "TikTok Followers", "Default", "6.00", 100, 50000, false, true, true, "tiktok_followers", "justanotherpanel"},

		// WhatsApp Services
		{40010000890098, "WhatsApp Group Members", "Default", "25.00", 50, 5000, false, true, true, "whatsapp_members", "justanotherpanel"},

		// Twitter/X Services
		{50010000890098, "Twitter Followers", "Default", "5.50", 100, 50000, false, true, true, "twitter_followers", "justanotherpanel"},
		{50020000890098, "Twitter Likes", "Default", "0.90", 50, 100000, false, true, true, "twitter_likes", "justanotherpanel"},
	}

	insertServiceQuery := `
		INSERT INTO smm_services 
		(service_id, name, type, rate, min_quantity, max_quantity, drip_feed, refill, cancel, category, provider_name, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, true)
		ON CONFLICT (service_id) DO NOTHING
	`

	for _, svc := range rawServices {
		_, err := conn.Exec(ctx, insertServiceQuery,
			svc.ServiceID, svc.Name, svc.Type, svc.Rate,
			svc.MinQty, svc.MaxQty, svc.DripFeed, svc.Refill, svc.Cancel,
			svc.Category, svc.ProviderName,
		)
		if err != nil {
			return fmt.Errorf("failed to seed smm_service %d: %w", svc.ServiceID, err)
		}
	}

	// ۲. درج نگاشت‌های curated (منوی کاربر)
	mappings := []struct {
		SmmServiceID int64
		Name         string
		Platform     string
		Category     string
		Description  string
		ButtonName   string
		SortOrder    int
		IsActive     bool
	}{
		// Telegram - Views
		{1001, "1K Post Views", "telegram", "views", "1000 بازدید پست تلگرام", "👁️ 1K بازدید", 1, true},
		{1002, "5K Post Views", "telegram", "views", "5000 بازدید پست تلگرام", "️ 5K بازدید", 2, true},
		{1003, "Story Views", "telegram", "views", "بازدید استوری تلگرام", "️ بازدید استوری", 3, true},

		// Telegram - Reactions
		{1004, "Positive Reaction", "telegram", "reactions", "ری‌اکشن مثبت (👍❤️🔥)", "❤️ ری‌اکشن مثبت", 1, true},

		// Telegram - Members
		{1005, "Real Members", "telegram", "members", "ممبر واقعی تلگرام", "👥 ممبر واقعی", 1, true},
		{1006, "Premium Members", "telegram", "members", "ممبر پرمیوم تلگرام", "💎 ممبر پرمیوم", 2, true},
		{1007, "Fake Members", "telegram", "members", "ممبر فیک ارزان", "👤 ممبر فیک", 3, true},

		// Telegram - Shares
		{1008, "Post Shares", "telegram", "shares", "فوروارد/شیر پست تلگرام", "🔄 فوروارد پست", 1, true},

		// Instagram - Likes
		{2001, "Real Likes", "instagram", "likes", "لایک واقعی اینستاگرام", "❤️ لایک واقعی", 1, true},
		{2002, "HQ Likes", "instagram", "likes", "لایک با کیفیت بالا", "💎 لایک HQ", 2, true},

		// Instagram - Followers
		{2003, "Real Followers", "instagram", "followers", "فالوور واقعی اینستاگرام", "👤 فالوور واقعی", 1, true},

		// Instagram - Views
		{2004, "Post Views", "instagram", "views", "بازدید پست اینستاگرام", "👁️ بازدید پست", 1, true},

		// TikTok
		{3001, "TikTok Views", "tiktok", "views", "بازدید تیک‌تاک", "👁️ بازدید", 1, true},
		{3002, "TikTok Likes", "tiktok", "likes", "لایک تیک‌تاک", "❤️ لایک", 2, true},
		{3003, "TikTok Followers", "tiktok", "followers", "فالوور تیک‌تاک", "👤 فالوور", 3, true},

		// WhatsApp
		{4001, "Group Members", "whatsapp", "members", "ممبر گروه واتس‌اپ", "👥 ممبر گروه", 1, true},

		// Twitter/X
		{5001, "Twitter Followers", "twitter", "followers", "فالوور توییتر", "👤 فالوور", 1, true},
		{5002, "Twitter Likes", "twitter", "likes", "لایک توییتر", "❤️ لایک", 2, true},
	}

	insertMappingQuery := `
		INSERT INTO smm_service_mappings 
		(smm_service_id, name, platform, category, description, button_name, sort_order, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT DO NOTHING
	`

	for _, m := range mappings {
		_, err := conn.Exec(ctx, insertMappingQuery,
			m.SmmServiceID, m.Name, m.Platform, m.Category,
			m.Description, m.ButtonName, m.SortOrder, m.IsActive,
		)
		if err != nil {
			return fmt.Errorf("failed to seed mapping for service %d: %w", m.SmmServiceID, err)
		}
	}

	return nil
}
