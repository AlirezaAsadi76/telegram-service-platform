package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

func main() {
	// ۱. اطلاعات کاربری که می‌خواهید لاگین را با آن تست کنید
	telegramUserID := int64(2053084840) // آیدی تلگرام خودتان یا یک آیدی تستی

	// ۲. بات توکن خود را از فایل کانفیگ اینجا کپی کنید (فقط برای تست محلی)
	botToken := "8873841525:AAGdvvgTW1hAzVjONNZ5apyXGmhEDVKbtZM"

	// ۳. ساخت داده‌های کاربر به فرمت JSON فشرده (بدون فاصله)
	userJSON := fmt.Sprintf(`{"id":%d}`, telegramUserID)
	authDate := fmt.Sprintf("%d", time.Now().Unix())

	// ۴. ساخت دیکشنری داده‌ها
	data := map[string]string{
		"user":      userJSON,
		"auth_date": authDate,
	}

	// مرتب‌سازی الفبایی کلیدها
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// ۲. ساخت رشته برای محاسبه HMAC (طبق استاندارد تلگرام با \n جدا می‌شوند)
	var dataCheckString []string
	for _, k := range keys {
		dataCheckString = append(dataCheckString, fmt.Sprintf("%s=%s", k, data[k]))
	}
	joinedDataForHmac := strings.Join(dataCheckString, "\n")

	// ۳. محاسبه هش
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(botToken))
	secretKey := h.Sum(nil)

	h = hmac.New(sha256.New, secretKey)
	h.Write([]byte(joinedDataForHmac))
	hash := hex.EncodeToString(h.Sum(nil))

	// ۴. ساخت رشته نهایی init_data (طبق استاندارد تلگرام با & جدا می‌شوند - تک خطی)
	var finalInitData []string
	for _, k := range keys {
		finalInitData = append(finalInitData, fmt.Sprintf("%s=%s", k, data[k]))
	}
	finalInitData = append(finalInitData, fmt.Sprintf("hash=%s", hash))

	result := strings.Join(finalInitData, "&")

	fmt.Println("✅ initData تولید شد (تک‌خطی و آماده کپی در Postman):")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(result)
	fmt.Println("--------------------------------------------------------------------------------")

}
