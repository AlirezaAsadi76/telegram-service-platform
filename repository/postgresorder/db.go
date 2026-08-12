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
	var providerID *uint64

	var metadata []byte
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.ProductType,
		&order.ProductID,
		&order.Quantity,
		&order.TargetLink,
		&order.Amount,
		&order.Currency,
		&order.Status,
		&order.ExternalOrderID,
		&providerID,
		&metadata,
		&order.CreatedAt,
		&order.UpdatedAt)
	
	order.ProviderID = providerID

	if len(metadata) > 0 {

		err = json.Unmarshal(
			metadata,
			&order.Metadata,
		)

	}
	return order, err

}
