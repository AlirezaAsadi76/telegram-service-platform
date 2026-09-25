package pricingservice

import (
	"errors"

	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/shopspring/decimal"
)

type SMMPricingRule struct {
	MinimumUnitPriceToman entity.Amount
	MinimumMultiplier     decimal.Decimal
	MaximumMultiplier     decimal.Decimal
	MultiplierScaleToman  entity.Amount
}

func (r SMMPricingRule) Validate() error {
	if r.MinimumUnitPriceToman.LessThan(entity.Amount(decimal.Zero)) ||
		r.MinimumUnitPriceToman.Equal(entity.Amount(decimal.Zero)) {
		return errors.New(
			"minimum SMM unit price must be greater than zero",
		)
	}

	if r.MinimumMultiplier.LessThan(decimal.Zero) ||
		r.MinimumMultiplier.Equal(decimal.Zero) {
		return errors.New(
			"minimum SMM multiplier must be greater than zero",
		)
	}

	if r.MaximumMultiplier.LessThan(r.MinimumMultiplier) {
		return errors.New(
			"maximum SMM multiplier must be greater than or equal to minimum multiplier",
		)
	}

	if r.MultiplierScaleToman.LessThan(entity.Amount(decimal.Zero)) ||
		r.MultiplierScaleToman.Equal(entity.Amount(decimal.Zero)) {
		return errors.New(
			"SMM multiplier scale must be greater than zero",
		)
	}

	return nil
}

func (r SMMPricingRule) CalculateSellingUnitPrice(baseCostPerThousand entity.Amount) (entity.Amount, error) {
	const op = "pricingservice.CalculateSMMSellingUnitPrice"

	if err := r.Validate(); err != nil {
		return entity.Amount{},
			richerror.New(op, err).
				WithKind(richerror.KindInvalid).
				WithCode(richerror.CodeInvalidInput).
				WithMessage(msgerror.InvalidPrice)
	}

	if baseCostPerThousand.LessThan(entity.Amount(decimal.Zero)) ||
		baseCostPerThousand.Equal(entity.Amount(decimal.Zero)) {
		return entity.Amount{},
			richerror.New(op, errors.New("base SMM cost must be greater than zero")).
				WithKind(richerror.KindInvalid).
				WithCode(richerror.CodeInvalidInput).
				WithMessage(msgerror.InvalidPrice)
	}

	base := baseCostPerThousand.Decimal()
	scale := r.MultiplierScaleToman.Decimal()

	multiplier := r.MinimumMultiplier.Add(
		r.MaximumMultiplier.Sub(r.MinimumMultiplier).Mul(
			scale.Div(base.Add(scale)),
		),
	)

	sellingPrice := base.Mul(multiplier)

	minimumPrice := r.MinimumUnitPriceToman.Decimal()

	if sellingPrice.LessThan(minimumPrice) {
		sellingPrice = minimumPrice
	}

	return entity.Amount(sellingPrice.Round(0)), nil
}
