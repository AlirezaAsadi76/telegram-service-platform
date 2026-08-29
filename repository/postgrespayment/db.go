package postgrespayment

import (
	"encoding/json"
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/repository/postgres"

	"github.com/shopspring/decimal"
)

type DB struct {
	Pool *postgres.DB
}

func New(pool *postgres.DB) *DB {
	return &DB{Pool: pool}
}

func scanPayment(row postgres.Scanner) (paymententity.Payment, error) {

	payment := paymententity.Payment{}
	var metadata []byte
	var amountStr string
	err := row.Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Method,
		&amountStr,
		&payment.Currency,
		&payment.Status,
		&payment.ExternalID,
		&payment.IdempotencyKey,
		&metadata,
		&payment.ExpiredAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if len(metadata) > 0 {

		err = json.Unmarshal(
			metadata,
			&payment.CallbackData,
		)

	}
	amount, sErr := decimal.NewFromString(amountStr)
	if sErr != nil {
		return payment, fmt.Errorf("failed to parse amount to decimal: %w", sErr)
	}
	payment.Amount = entity.Amount(amount)
	return payment, err

}
