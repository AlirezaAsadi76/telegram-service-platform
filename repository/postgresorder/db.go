package postgresorder

import (
	"encoding/json"
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/repository/postgres"

	"github.com/jackc/pgx/v5/pgtype"
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
	var amountNumeric pgtype.Numeric
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.ProductType,
		&order.ProductID,
		&order.Quantity,
		&order.TargetLink,
		&amountNumeric,
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

	if amountNumeric.Valid {
		amountDecimal, err := numericToDecimal(amountNumeric)
		if err != nil {
			return order, fmt.Errorf("failed to convert numeric to decimal: %w", err)
		}
		order.Amount = entity.Amount(amountDecimal)
	} else {
		order.Amount = entity.Amount(decimal.Zero)
	}

	return order, err

}

func numericToDecimal(n pgtype.Numeric) (decimal.Decimal, error) {
	if !n.Valid {
		return decimal.Zero, nil
	}

	if n.Int != nil {

		intStr := n.Int.String()
		d, err := decimal.NewFromString(intStr)
		if err != nil {
			return decimal.Zero, err
		}
		if n.Exp != 0 {
			d = d.Shift(n.Exp)
		}
		return d, nil
	}

	return decimal.Zero, nil
}
