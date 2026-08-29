package postgresorder

import (
	"encoding/json"
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/repository/postgres"

	"github.com/shopspring/decimal"
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
	var amountStr string
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.ProductType,
		&order.ProductID,
		&order.Quantity,
		&order.TargetLink,
		&amountStr,
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

	amount, sErr := decimal.NewFromString(amountStr)
	if sErr != nil {
		return order, fmt.Errorf("failed to parse amount to decimal: %w", sErr)
	}
	order.Amount = entity.Amount(amount)

	return order, err

}
