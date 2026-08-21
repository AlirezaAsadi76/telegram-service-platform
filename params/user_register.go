package params

import (
	"telegram-service-platform/entity"
)

type GetOrRegisterRequest struct {
	TelegramID int64
	Username   string
	FirstName  string
	LastName   string
	Role       entity.Role
}
type GetOrRegisterResponse struct {
	UserInfo UserInfo `json:"user_info"`
	IsNew    bool     `json:"is_new"`
}
