package postgresnotification

import (
	"context"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"
)

func (db *DB) UpdateStatus(ctx context.Context, id uint64, status notificationentity.NotificationStatus) error {
	const Op = "postgres_notification.UpdateStatus"

	query := `UPDATE notifications SET status = $1, sent_at = $2 WHERE id = $3`
	var sentAt *time.Time
	if status == notificationentity.NotificationStatusSent {
		sentAt = new(time.Now())
	}
	_, err := db.Pool.Connection().Exec(ctx, query, status, sentAt, id)
	if err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}

	return nil
}
