package postgrescheckouttesting

import (
	"context"
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/repository/postgrescheckout"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

func getRefundTransactionCount(t *testing.T, pool *pgxpool.Pool, orderID uint64) int {
	t.Helper()

	idempotencyKey := fmt.Sprintf(
		postgrescheckout.IdempotencyRefund,
		orderID,
	)

	var count int

	err := pool.QueryRow(
		context.Background(),
		`
		SELECT COUNT(*)
		FROM wallet_transactions
		WHERE idempotency_key = $1
		`,
		idempotencyKey,
	).Scan(&count)

	if err != nil {
		t.Fatalf(
			"get refund transaction count: %v",
			err,
		)
	}

	return count
}

type RefundTransaction struct {
	txID            uint64
	txType          string
	amount          entity.Amount
	status          string
	referenceID     string
	idempotency_key string
}

func getRefundTransaction(t *testing.T, pool *pgxpool.Pool, orderID uint64) RefundTransaction {
	t.Helper()

	idempotencyKey := fmt.Sprintf(
		postgrescheckout.IdempotencyRefund,
		orderID,
	)

	var result RefundTransaction
	var amount string
	err := pool.QueryRow(
		context.Background(),
		`
		SELECT
			id,
			type::text,
			amount::text,
			status::text,
			reference_id,
			idempotency_key
		FROM wallet_transactions
		WHERE idempotency_key = $1
		`,
		idempotencyKey,
	).Scan(
		&result.txID,
		&result.txType,
		&amount,
		&result.status,
		&result.referenceID,
		&result.idempotency_key,
	)

	if err != nil {
		t.Fatalf(
			"get refund transaction: %v",
			err,
		)
	}

	value, err := decimal.NewFromString(amount)
	if err != nil {
		t.Fatalf(
			"parse refund amount %q: %v",
			amount,
			err,
		)
	}
	result.amount = entity.Amount(value)

	return result
}
