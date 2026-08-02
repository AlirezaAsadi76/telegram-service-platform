package params

import "telegram-service-platform/service/entity"

type UserInfo struct {
	Id         uint64      `json:"id"`
	TelegramID int64       `json:"telegram_id"`
	Username   string      `json:"username"`
	Role       entity.Role `json:"role"`
}
