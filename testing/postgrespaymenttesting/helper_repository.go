package postgrespaymenttesting

import (
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgrespayment"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newPaymentRepository(t *testing.T, pool *pgxpool.Pool) *postgrespayment.DB {
	t.Helper()
	transactionProvider := postgres.NewTransactionProvider(pool)
	return postgrespayment.NewWithExecutor(
		pool,
		transactionProvider,
	)
}
