package userparams

import (
	"telegram-service-platform/entity"
)

type GetOrRegisterRequest struct {
	TelegramID int64       `json:"telegram_id"`
	Username   string      `json:"username"`
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	Role       entity.Role `json:"role"`
}
type GetOrRegisterResponse struct {
	UserInfo UserInfo `json:"user_info"`
	IsNew    bool     `json:"is_new"`
}
