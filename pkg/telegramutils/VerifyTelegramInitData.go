package telegramutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func VerifyTelegramInitData(initData, botToken string) (url.Values, error) {

	parsedData, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("invalid init data format: %w", err)
	}

	clientHash := parsedData.Get("hash")
	parsedData.Del("hash")

	var keys []string
	for k := range parsedData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// ۴. ساخت رشته‌ی data_check_string به فرمت: key=value\nkey=value
	var dataCheckString []string
	for _, k := range keys {
		dataCheckString = append(dataCheckString, fmt.Sprintf("%s=%s", k, parsedData.Get(k)))
	}
	joinedData := strings.Join(dataCheckString, "\n")

	// ۵. ساخت کلید مخفی (Secret Key) با استفاده از HMAC-SHA256 روی رشته‌ی "WebAppData" و Bot Token
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(botToken))
	secretKey := h.Sum(nil)

	// ۶. محاسبه‌ی هش داده‌های مرتب‌شده با استفاده از کلید مخفی ساخته‌شده
	h = hmac.New(sha256.New, secretKey)
	h.Write([]byte(joinedData))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	// ۷. مقایسه‌ی امن (ضد حملات Timing Attack) هش محاسبه‌شده با هش ارسالی از کلاینت
	if !hmac.Equal([]byte(calculatedHash), []byte(clientHash)) {
		return nil, fmt.Errorf("invalid telegram signature")
	}

	return parsedData, nil
}
