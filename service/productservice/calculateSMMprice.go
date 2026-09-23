package productservice

import (
	"context"
	"errors"

	"telegram-service-platform/entity"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"github.com/shopspring/decimal"
)

func (s Service) CalculateSMMPrice(ctx context.Context, req productparams.CalculateSMMPriceRequest) (productparams.CalculateSMMPriceResponse, error) {
	const op = "productservice.CalculateSMMPrice"

	if req.MappingID <= 0 || req.Quantity <= 0 {
		return productparams.CalculateSMMPriceResponse{},
			richerror.New(op, errors.New("mapping ID and quantity must be positive")).
				WithKind(richerror.KindValidation).
				WithCode(richerror.CodeInvalidInput).
				WithMessage(msgerror.InvalidInput)
	}

	mapping, err := s.repository.SMMMappingGetByID(ctx, req.MappingID)
	if err != nil {
		return productparams.CalculateSMMPriceResponse{},
			richerror.New(op, err).
				WithKind(richerror.KindNotFound).
				WithCode(richerror.CodeProductNotFound).
				WithMessage(msgerror.ProductNotFound)
	}

	if !mapping.IsActive {
		return productparams.CalculateSMMPriceResponse{},
			richerror.New(
				op,
				errors.New("smm mapping is inactive"),
			).
				WithKind(richerror.KindNotFound).
				WithCode(richerror.CodeProductNotFound).
				WithMessage(msgerror.ProductNotFound)
	}

	service, smmErr := s.repository.SMMServiceGetByD(ctx, mapping.SmmServiceId)
	if smmErr != nil {
		return productparams.CalculateSMMPriceResponse{},
			richerror.New(op, smmErr).
				WithKind(richerror.KindNotFound).
				WithCode(richerror.CodeProductNotFound).
				WithMessage(msgerror.ProductNotFound)
	}

	if !service.IsActive {
		return productparams.CalculateSMMPriceResponse{},
			richerror.New(op, errors.New("smm service is inactive")).
				WithKind(richerror.KindNotFound).
				WithCode(richerror.CodeProductNotFound).
				WithMessage(msgerror.ProductNotFound)
	}

	if service.Rate.LessThan(
		entity.Amount(decimal.NewFromInt(0)),
	) || service.Rate.Equal(
		entity.Amount(decimal.NewFromInt(0)),
	) {
		return productparams.CalculateSMMPriceResponse{},
			richerror.New(op, errors.New("smm service rate must be greater than zero")).
				WithKind(richerror.KindInvalid).
				WithCode(richerror.CodeInvalidInput).
				WithMessage(msgerror.InvalidPrice)
	}

	qty := entity.Amount(decimal.NewFromInt(req.Quantity))

	usdPrice := service.Rate.Mul(qty).Div(entity.Amount(decimal.NewFromInt(1000)))

	price, cErr := s.pricingSVc.CalculatePrice(ctx, usdPrice)
	if cErr != nil {
		return productparams.CalculateSMMPriceResponse{},
			richerror.New(op, cErr).
				WithKind(richerror.KindUnexpected).
				WithMessage(msgerror.Unexpected)
	}

	return productparams.CalculateSMMPriceResponse{
		MappingID: req.MappingID,
		ServiceID: service.Id,
		Rate:      service.Rate,
		Price:     price,
	}, nil
}
