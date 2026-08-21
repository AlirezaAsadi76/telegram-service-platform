package postgresuser

import (
	"context"
	"database/sql"
	"errors"
	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/repository/postgres"
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
			last_seen_at,
			created_at,
			updated_at
		FROM users
		WHERE telegram_id = $1`

	row := db.Pool.Connection().QueryRow(ctx, query, telegramID)
	userCreated, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, richerror.New(Op, err).
				WithKind(richerror.KindNotFound).
				WithMessage(msgerror.UserNotFound)
		}
		return nil, richerror.New(Op, err).WithKind(richerror.KindUnexpected).WithMessage(msgerror.Unexpected)
	}
	return &userCreated, nil
}

func scanUser(row postgres.Scanner) (entity.User, error) {
	user := entity.User{}
	err := row.Scan(&user.ID, &user.TelegramID, &user.Username, &user.FirstName,
		&user.LastName, &user.Role, &user.LastSeenAt, &user.CreatedAt, &user.UpdatedAt)

	return user, err

}
