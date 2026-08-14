package notificationparams

import "telegram-service-platform/entity/notificationentity"

type UpdateStatusRequest struct {
	Id     uint64
	Status notificationentity.NotificationStatus
}
