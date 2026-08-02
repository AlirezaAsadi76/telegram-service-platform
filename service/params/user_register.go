package params

import "telegram-service-platform/service/entity"

type RegisterRequest struct {
	TelegramID int64
	Username   string
	FirstName  string
	LastName   string
	Role       entity.Role
}
type RegisterResponse struct {
	UserInfo UserInfo `json:"user_info"`
}
