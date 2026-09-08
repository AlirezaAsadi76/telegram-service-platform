package postgrespayment

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgUniqueViolationCode = "23505"

	paymentIdempotencyConstraint = "ux_payments_idempotency_key"
	paymentActiveOrderConstraint = "ux_payments_active_order"
)

func getUniqueViolationConstraint(err error) (string, bool) {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return "", false
	}

	if pgErr.Code != pgUniqueViolationCode {
		return "", false
	}

	return pgErr.ConstraintName, true
}
