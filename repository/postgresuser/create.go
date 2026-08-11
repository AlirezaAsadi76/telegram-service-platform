package postgresuser

import (
	"context"
	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (db *DB) Create(ctx context.Context, user *entity.User) error {
	const Op = "postgresuser.create"
	// TODO - create wallet for user when create user
	err := db.Pool.Connection().QueryRow(ctx, `
							INSERT INTO users (telegram_id, first_name, last_name, username, role) 
							VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`,
		user.TelegramID, user.FirstName, user.LastName, user.Username, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return richerror.New(Op, err).
			WithKind(richerror.KindUnexpected).
			WithMessage(msgerror.Unexpected)
	}
	return nil
}
