package notificationparams

import "telegram-service-platform/entity/notificationentity"

type CreateRequest struct {
	UserID  uint64
	Type    notificationentity.NotificationType
	Payload map[string]any
}
