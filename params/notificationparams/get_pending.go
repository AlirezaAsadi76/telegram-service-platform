package notificationparams

import "telegram-service-platform/entity/notificationentity"

type GetPendingRequest struct {
	Limit int
}

type GetPendingResponse struct {
	Notifications []notificationentity.Notification
}
