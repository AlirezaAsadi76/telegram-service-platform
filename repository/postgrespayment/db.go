package postgrespayment

import (
	"encoding/json"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/repository/postgres"
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
	err := row.Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Method,
		&payment.Amount,
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
	return payment, err

}
