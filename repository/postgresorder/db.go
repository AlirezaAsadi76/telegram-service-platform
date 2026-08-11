package postgresorder

import (
	"encoding/json"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/repository/postgres"
)

type DB struct {
	Pool *postgres.DB
}

func New(pool *postgres.DB) *DB {

	return &DB{Pool: pool}
}

func scanOrder(row postgres.Scanner) (orderentity.Order, error) {
	order := orderentity.Order{}
	var metadata []byte
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.ProductType,
		&order.ProductID,
		&order.Quantity,
		&order.Amount,
		&order.Currency,
		&order.Status,
		&metadata,
		&order.CreatedAt,
		&order.UpdatedAt)

	if len(metadata) > 0 {

		err = json.Unmarshal(
			metadata,
			&order.Metadata,
		)

	}
	return order, err

}
