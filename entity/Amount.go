package entity

import (
	"database/sql/driver"

	"github.com/shopspring/decimal"
)

type Amount decimal.Decimal

func (a Amount) Value() (driver.Value, error) {
	return decimal.Decimal(a), nil
}
