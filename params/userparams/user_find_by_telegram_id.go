package userparams

import "telegram-service-platform/entity"

type FindUserByTelegramIDRequest struct {
	TelegramID entity.TelegramId
}

type FindUserByTelegramIDResponse struct {
	UserInfo UserInfo
	Found    bool
}
