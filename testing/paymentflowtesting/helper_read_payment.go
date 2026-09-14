package paymentflowtesting

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/repository/postgres"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

func readPayment(t *testing.T, pool *pgxpool.Pool, paymentID uint64) paymententity.Payment {
	t.Helper()

	row := pool.QueryRow(
		context.Background(),
		`
			SELECT
				id,
				order_id,
				user_id,
				method,
				amount,
				currency,
				status,
				external_id,
				provider_reference_id,
				payment_url,
				idempotency_key,
				callback_data,
				expired_at,
				created_at,
				updated_at
			FROM payments
			WHERE id = $1
		`,
		paymentID,
	)

	payment, err := scanPayment(row)

	if err != nil {
		t.Fatalf("read persisted payment: %v", err)
	}
	return payment
}

func scanPayment(row postgres.Scanner) (paymententity.Payment, error) {

	payment := paymententity.Payment{}
	var metadata []byte
	var amountStr string
	var providerRef, paymentUrl, idempotencyKey sql.NullString
	var expiresAt sql.NullTime
	err := row.Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Method,
		&amountStr,
		&payment.Currency,
		&payment.Status,
		&payment.ExternalID,
		&providerRef,
		&paymentUrl,
		&idempotencyKey,
		&metadata,
		&expiresAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if len(metadata) > 0 {

		err = json.Unmarshal(
			metadata,
			&payment.CallbackData,
		)

	}
	payment.ProviderReferenceID = providerRef.String
	payment.PaymentURL = paymentUrl.String
	payment.IdempotencyKey = idempotencyKey.String
	payment.ExpiredAt = expiresAt.Time
	if amountStr != "" {
		amount, sErr := decimal.NewFromString(amountStr)
		if sErr != nil {
			return payment, fmt.Errorf("failed to parse amount to decimal: %w", sErr)
		}
		payment.Amount = entity.Amount(amount)
	}
	return payment, err

}
