package authservice

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	TelegramId int64  `json:"telegram_id,omitempty"`
	UserId     int64  `json:"user_id,omitempty"`
	Role       string `json:"role,omitempty"`
}
