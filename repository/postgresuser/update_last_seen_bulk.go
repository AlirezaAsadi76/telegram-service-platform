package postgresuser

import (
	"context"
	"fmt"
	"strings"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) UpdateLastSeenBulk(ctx context.Context, telegramIDs []int64) (int64, error) {
	const op = "postgresuser.UpdateLastSeenBulk"

	if len(telegramIDs) == 0 {
		return 0, nil
	}

	placeholders := make([]string, len(telegramIDs))
	args := make([]interface{}, len(telegramIDs))
	for i, id := range telegramIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		UPDATE users
		SET last_seen_at = NOW(), updated_at = NOW()
		WHERE telegram_id IN (%s)
	`, strings.Join(placeholders, ","))

	result, err := db.Pool.Connection().Exec(ctx, query, args...)
	if err != nil {
		return 0, richerror.New(op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.QueryFailed)
	}

	rowsAffected := result.RowsAffected()

	return rowsAffected, nil
}
