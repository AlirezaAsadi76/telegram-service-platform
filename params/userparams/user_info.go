package userparams

import (
	"telegram-service-platform/entity"
	"time"
)

type UserInfo struct {
	Id         uint64      `json:"id"`
	TelegramID int64       `json:"telegram_id"`
	Username   string      `json:"username"`
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	Role       entity.Role `json:"role"`
	CreatedAt  time.Time   `json:"created_at"`
}
