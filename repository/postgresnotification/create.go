package postgresnotification

import (
	"context"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) Create(ctx context.Context, notification *notificationentity.Notification) error {
	const Op = "postgres_notification.Create"
	query := `
		INSERT INTO notifications (user_id, type, status, payload, retry_count)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := db.Pool.Connection().QueryRow(ctx, query,
		notification.UserID, notification.Type, notification.Status, notification.Payload, notification.RetryCount,
	).Scan(&notification.ID, &notification.CreatedAt)

	if err != nil {
		return richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryScanFailed)
	}
	return nil
}
