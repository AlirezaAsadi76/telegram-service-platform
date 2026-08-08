package postgresuser

import (
	"context"
	"database/sql"
	"errors"
	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) FindUserByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error) {
	const Op = "postgresuser.find_user_by_telegram_id"

	query := `SELECT
			id,
			telegram_id,
			username,
			first_name,
			last_name,
			role,
			created_at,
			updated_at
		FROM users
		WHERE telegram_id = $1`
	var user entity.User
	err := db.Pool.Connection().QueryRow(ctx, query, telegramID).
		Scan(&user.ID, &user.TelegramID, &user.Username, &user.FirstName,
			&user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, richerror.New(Op, err).
				WithKind(richerror.KindNotFound).
				WithMessage(msgerror.UserNotFound)
		}
		return nil, richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}
	return &user, nil
}
