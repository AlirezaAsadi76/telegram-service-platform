package postgresnotification

import (
	"encoding/json"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/repository/postgres"
)

type DB struct {
	Pool *postgres.DB
}

func New(pool *postgres.DB) *DB {

	return &DB{Pool: pool}
}

func scanNotification(row postgres.Scanner) (notificationentity.Notification, error) {

	notification := notificationentity.Notification{}
	var payload []byte

	err := row.Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Type,
		&notification.Status,
		&payload,
		&notification.RetryCount,
		&notification.CreatedAt,
		&notification.SentAt,
	)

	if len(payload) > 0 {

		err = json.Unmarshal(
			payload,
			&notification.Payload,
		)

	}
	return notification, err

}
