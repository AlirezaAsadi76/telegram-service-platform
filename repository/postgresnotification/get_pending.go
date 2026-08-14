package postgresnotification

import (
	"context"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) GetPending(ctx context.Context, limit int) ([]notificationentity.Notification, error) {
	const Op = "postgresnotification.GetPending"

	query := `
		SELECT id, user_id, type, status, payload, retry_count, created_at, sent_at
		FROM notifications
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
		LIMIT $1
	`
	rows, qErr := db.Pool.Connection().Query(ctx, query, limit)
	if qErr != nil {
		return nil, richerror.New(Op, qErr).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryScanFailed)
	}

	defer rows.Close()

	notifications := make([]notificationentity.Notification, 0)

	for rows.Next() {

		notif, sErr := scanNotification(rows)
		if sErr != nil {
			return nil, richerror.New(Op, sErr).WithKind(richerror.KindScanFailure).WithMessage(msgerror.QueryScanFailed)
		}

		notifications = append(
			notifications,
			notif,
		)
	}

	return notifications, nil
}
