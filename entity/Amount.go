package entity

import (
	"github.com/shopspring/decimal"
)

type Amount decimal.Decimal

func (a *Amount) UnmarshalJSON(data []byte) error {

	return (*decimal.Decimal)(a).UnmarshalJSON(data)
}

func (a Amount) MarshalJSON() ([]byte, error) {

	return (decimal.Decimal)(a).MarshalJSON()
}

func (a Amount) GreaterThan(b Amount) bool {
	return decimal.Decimal(a).GreaterThan(decimal.Decimal(b))
}

func (a Amount) Mul(b Amount) Amount {

	return Amount(decimal.Decimal(a).Mul(decimal.Decimal(b)))
}

func (a Amount) Round(b int32) Amount {
	return Amount(decimal.Decimal(a).Round(b))
}

func (a Amount) Div(b Amount) Amount {
	return Amount(decimal.Decimal(a).Div(decimal.Decimal(b)))
}
func (a Amount) Sub(b Amount) Amount {
	return Amount(decimal.Decimal(a).Sub(decimal.Decimal(b)))
}
func (a Amount) Add(b Amount) Amount {
	return Amount(decimal.Decimal(a).Add(decimal.Decimal(b)))
}

func (a Amount) MulInt(b int64) Amount {
	qty := decimal.NewFromInt(b)
	return Amount(decimal.Decimal(a).Mul(qty))
}
func (a Amount) LessThan(b Amount) bool {
	return decimal.Decimal(a).LessThan(decimal.Decimal(b))
}
func (a Amount) String() string {
	return decimal.Decimal(a).String()
}
