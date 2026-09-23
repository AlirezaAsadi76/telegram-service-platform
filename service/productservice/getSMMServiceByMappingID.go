package productservice

import (
	"context"
	"errors"
	"time"

	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	"go.uber.org/zap"
)

func (s Service) GetSMMServiceByMappingID(ctx context.Context, mappingID int64) (*smmentity.SMM, error) {
	const op = "productservice.GetSMMServiceByMappingID"

	start := time.Now()

	if mappingID <= 0 {
		return nil,
			richerror.New(op, errors.New("mapping ID must be greater than zero")).
				WithKind(richerror.KindValidation).
				WithCode(richerror.CodeInvalidInput).
				WithMessage(msgerror.InvalidInput)
	}

	mappingResponse, err := s.GetSMMMappingByID(ctx,
		productparams.GetSmmMappingByIDRequest{
			Id: mappingID,
		},
	)
	if err != nil {
		logger.Logger.Error(
			"failed to resolve SMM mapping",
			zap.String("op", op),
			zap.Int64("mapping_id", mappingID),
			zap.Error(err),
			zap.Duration(
				"duration",
				time.Since(start),
			),
		)

		return nil, richerror.New(op, err)
	}

	if mappingResponse.SmmMapping == nil {
		return nil,
			richerror.New(op, errors.New("resolved SMM mapping is nil")).
				WithKind(richerror.KindNotFound).
				WithCode(richerror.CodeProductNotFound).
				WithMessage(msgerror.ProductNotFound)
	}

	smmResponse, smmErr := s.GetSMMServiceByID(
		ctx,
		productparams.GetSmmServiceByIDRequest{
			Id: mappingResponse.SmmMapping.SmmServiceId,
		},
	)
	if smmErr != nil {
		logger.Logger.Error(
			"failed to resolve SMM service from mapping",
			zap.String("op", op),
			zap.Int64("mapping_id", mappingID),
			zap.Int64("smm_id", mappingResponse.SmmMapping.SmmServiceId),
			zap.Error(smmErr),
			zap.Duration("duration", time.Since(start)),
		)

		return nil, richerror.New(op, smmErr)
	}

	if smmResponse.Smm == nil {
		return nil,
			richerror.New(
				op,
				errors.New("resolved SMM service is nil"),
			).
				WithKind(richerror.KindNotFound).
				WithCode(richerror.CodeProductNotFound).
				WithMessage(msgerror.ProductNotFound)
	}

	logger.Logger.Debug(
		"SMM service resolved from mapping",
		zap.String("op", op),
		zap.Int64("mapping_id", mappingID),
		zap.Int64("smm_id", smmResponse.Smm.Id),
		zap.Int64("provider_service_id", smmResponse.Smm.Service),
		zap.String("provider_name", smmResponse.Smm.ProviderName),
		zap.Duration("duration", time.Since(start)),
	)

	return smmResponse.Smm, nil
}
