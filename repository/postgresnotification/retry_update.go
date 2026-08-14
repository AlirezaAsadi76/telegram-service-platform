package postgresnotification

import (
	"context"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) UpdateRetryCount(ctx context.Context, id uint64, retryCount int) error {
	const Op = "postgres_notification.UpdateRetryCount"

	query := `UPDATE notifications SET retry_count = $1 WHERE id = $2`
	_, err := db.Pool.Connection().Exec(ctx, query, retryCount, id)

	if err != nil {
		return richerror.New(Op, err).WithKind(richerror.KindQueryFailure).WithMessage(msgerror.QueryFailed)
	}
	return err
}
