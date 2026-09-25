package pricingservice

import (
	"context"
	"errors"

	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/shopspring/decimal"
)

func (s Service) CalculateSMMSellingPrice(ctx context.Context, ratePerThousand entity.Amount, quantity int64) (productentity.Price, error) {
	const op = "pricingservice.CalculateSMMSellingPrice"

	if quantity <= 0 {
		return productentity.Price{},
			richerror.New(op, errors.New("quantity must be greater than zero")).
				WithKind(richerror.KindInvalid).
				WithCode(richerror.CodeInvalidInput).
				WithMessage(msgerror.InvalidInput)
	}

	if ratePerThousand.LessThan(entity.Amount(decimal.Zero)) ||
		ratePerThousand.Equal(entity.Amount(decimal.Zero)) {
		return productentity.Price{},
			richerror.New(op, errors.New("SMM rate must be greater than zero")).
				WithKind(richerror.KindInvalid).
				WithCode(richerror.CodeInvalidInput).
				WithMessage(msgerror.InvalidPrice)
	}

	basePrice, err := s.CalculatePrice(ctx, ratePerThousand.MulInt(quantity).Div(entity.Amount(decimal.NewFromInt(1000))))
	if err != nil {
		return productentity.Price{}, richerror.New(op, err)
	}

	baseUnitToman := basePrice.Toman.Mul(entity.Amount(decimal.NewFromInt(1000))).Div(entity.Amount(decimal.NewFromInt(quantity)))

	rule, rErr := s.smmPricingRule()
	if rErr != nil {
		return productentity.Price{}, richerror.New(op, rErr)
	}

	sellingUnitToman, pErr := rule.CalculateSellingUnitPrice(baseUnitToman)
	if pErr != nil {
		return productentity.Price{}, richerror.New(op, pErr)
	}

	totalToman := sellingUnitToman.
		MulInt(quantity).
		Div(entity.Amount(decimal.NewFromInt(1000))).
		Round(0)

	if totalToman.Equal(entity.Amount(decimal.Zero)) {
		return productentity.Price{},
			richerror.New(op, errors.New("calculated SMM selling price is zero")).
				WithKind(richerror.KindInvalid).
				WithCode(richerror.CodeInvalidInput).
				WithMessage(msgerror.InvalidPrice)
	}

	baseToman := basePrice.Toman.Decimal()
	baseUSD := basePrice.USD.Decimal()

	if baseToman.Equal(decimal.Zero) ||
		baseUSD.Equal(decimal.Zero) {
		return productentity.Price{},
			richerror.New(op, errors.New("invalid base SMM price")).
				WithKind(richerror.KindInvalid).
				WithCode(richerror.CodeInvalidInput).
				WithMessage(msgerror.InvalidPrice)
	}

	exchangeRateUSDToToman := baseToman.Div(baseUSD)

	finalUSD := totalToman.Decimal().
		Div(exchangeRateUSDToToman)

	tonRateUSD := basePrice.USD.Decimal().
		Div(basePrice.TON.Decimal())

	if tonRateUSD.Equal(decimal.Zero) {
		return productentity.Price{},
			richerror.New(
				op,
				errors.New("invalid USD/TON conversion rate"),
			).
				WithKind(richerror.KindInvalid).
				WithCode(richerror.CodeInvalidInput).
				WithMessage(msgerror.InvalidPrice)
	}

	finalTON := finalUSD.Div(tonRateUSD)

	return productentity.Price{
		USD:   entity.Amount(finalUSD),
		USDT:  entity.Amount(finalUSD),
		TON:   entity.Amount(finalTON),
		Toman: totalToman,
	}, nil
}

func (s Service) smmPricingRule() (SMMPricingRule, error) {
	const op = "pricingservice.smmPricingRule"

	minimumPrice, err := decimal.NewFromString(
		s.config.SMMMinimumUnitPriceToman,
	)
	if err != nil {
		return SMMPricingRule{}, richerror.New(op, err).
			WithKind(richerror.KindInvalid).
			WithCode(richerror.CodeInvalidInput).
			WithMessage(msgerror.InvalidPrice)
	}

	minimumMultiplier, err := decimal.NewFromString(
		s.config.SMMMinimumMultiplier,
	)
	if err != nil {
		return SMMPricingRule{}, richerror.New(op, err).
			WithKind(richerror.KindInvalid).
			WithCode(richerror.CodeInvalidInput).
			WithMessage(msgerror.InvalidPrice)
	}

	maximumMultiplier, err := decimal.NewFromString(
		s.config.SMMMaximumMultiplier,
	)
	if err != nil {
		return SMMPricingRule{}, richerror.New(op, err).
			WithKind(richerror.KindInvalid).
			WithCode(richerror.CodeInvalidInput).
			WithMessage(msgerror.InvalidPrice)
	}

	scale, err := decimal.NewFromString(
		s.config.SMMMultiplierScaleToman,
	)
	if err != nil {
		return SMMPricingRule{}, richerror.New(op, err).
			WithKind(richerror.KindInvalid).
			WithCode(richerror.CodeInvalidInput).
			WithMessage(msgerror.InvalidPrice)
	}

	return SMMPricingRule{
		MinimumUnitPriceToman: entity.Amount(minimumPrice),
		MinimumMultiplier:     minimumMultiplier,
		MaximumMultiplier:     maximumMultiplier,
		MultiplierScaleToman:  entity.Amount(scale),
	}, nil
}
