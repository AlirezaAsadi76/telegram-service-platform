package params

import "telegram-service-platform/entity"

type TrackUserActivityRequest struct {
	TelegramID entity.TelegramId
}

type TrackUserActivityResponse struct {
	Tracked bool
}

type GetActiveUsersRequest struct{}

type GetActiveUsersResponse struct {
	ActiveTelegramIDs []entity.TelegramId
	Count             int
}

type SyncLastSeenRequest struct{}

type SyncLastSeenResponse struct {
	Synced entity.TelegramId
}
