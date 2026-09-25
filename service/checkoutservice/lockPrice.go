package checkoutservice

import (
	"context"
	"errors"
	"telegram-service-platform/entity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"time"

	"go.uber.org/zap"
)

func (s *Service) LockSMMPrice(ctx context.Context, req checkoutparams.LockProductPriceRequest) (*checkoutparams.LockProductPriceResponse, error) {
	const op = "checkoutservice.LockSMMPrice"

	if s.productSvc == nil {
		return nil, richerror.New(op, errors.New("product price resolver is unavailable")).
			WithKind(richerror.KindDependencyFailure).
			WithCode(richerror.CodeSMMProviderUnavailable).
			WithMessage(msgerror.ExternalServiceFailed)
	}

	ttl, err := s.priceLockTTL(req.PaymentMethod)
	if err != nil {
		return nil, richerror.New(op, err).
			WithKind(richerror.KindValidation).
			WithCode(richerror.CodeInvalidInput).
			WithMessage(msgerror.InvalidInput)
	}

	priceResp, cErr := s.productSvc.CalculateSMMPrice(
		ctx,
		productparams.CalculateSMMPriceRequest{
			MappingID: int64(req.ProductID),
			Quantity:  req.Quantity,
		},
	)
	if cErr != nil {
		return nil, richerror.New(op, cErr)
	}

	amount, pErr := priceForCurrency(
		priceResp,
		req.Currency,
	)
	if pErr != nil {
		return nil, richerror.New(op, pErr).
			WithKind(richerror.KindInvalid).
			WithCode(richerror.CodeInvalidInput).
			WithMessage(msgerror.InvalidPrice)
	}

	if amount.Equal(entity.Amount{}) ||
		amount.LessThan(entity.Amount{}) {

		return nil, richerror.New(
			op,
			errors.New("locked SMM price must be greater than zero"),
		).
			WithKind(richerror.KindInvalid).
			WithCode(richerror.CodeInvalidInput).
			WithMessage(msgerror.InvalidPrice)
	}

	now := time.Now()
	expiresAt := now.Add(ttl)

	metrics.CheckoutPriceLocks.WithLabelValues(string(req.PaymentMethod), "created").Inc()

	logger.Logger.Info(
		"SMM price locked",
		zap.Uint64("user_id", req.UserID),
		zap.Uint64("mapping_id", req.ProductID),
		zap.Int64("quantity", req.Quantity),
		zap.String("payment_method", string(req.PaymentMethod)),
		zap.String("currency", string(req.Currency)),
		zap.String("amount", amount.String()),
		zap.Time("expires_at", expiresAt),
	)

	return &checkoutparams.LockProductPriceResponse{
		Amount:        amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		LockedAt:      now.Unix(),
		ExpiresAt:     expiresAt.Unix(),
	}, nil
}

func (s *Service) priceLockTTL(method entity.PriceLockPaymentMethod) (time.Duration, error) {
	switch method {
	case entity.PriceLockPaymentMethodWallet:
		if s.config.WalletPriceLockTTL > 0 {
			return s.config.WalletPriceLockTTL, nil
		}

	case entity.PriceLockPaymentMethodZarinpal:
		if s.config.ZarinpalPriceLockTTL > 0 {
			return s.config.ZarinpalPriceLockTTL, nil
		}

	case entity.PriceLockPaymentMethodCrypto:
		if s.config.CryptoPriceLockTTL > 0 {
			return s.config.CryptoPriceLockTTL, nil
		}
	}

	return 0, errors.New("price lock TTL is not configured")
}

func priceForCurrency(response productparams.CalculateSMMPriceResponse, currency entity.Currency) (entity.Amount, error) {
	switch currency {
	case entity.CurrencyTOMAN:
		return response.Price.Toman, nil

	case entity.CurrencyUSD:
		return response.Price.USD, nil

	case entity.CurrencyUSDT:
		return response.Price.USDT, nil

	case entity.CurrencyTON:
		return response.Price.TON, nil

	default:
		return entity.Amount{},
			errors.New("unsupported SMM price currency")
	}
}
