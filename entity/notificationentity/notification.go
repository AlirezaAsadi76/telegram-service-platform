package notificationentity

import "time"

type Notification struct {
	ID         uint64
	UserID     uint64
	Type       NotificationType
	Status     NotificationStatus
	Payload    map[string]any
	RetryCount int
	CreatedAt  time.Time
	SentAt     *time.Time
}
