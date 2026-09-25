package checkoutValidator

import (
	"errors"
	"telegram-service-platform/entity"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/pkg/wrapper"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v *Validator) ValidatePriceLock(req checkoutparams.WalletPurchaseRequest, method entity.PriceLockPaymentMethod) error {

	const op = "checkout.ValidatePriceLock"

	vErr := validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required, validation.Min(1)),
		validation.Field(&req.Amount, validation.Required),
	)

	if vErr != nil {
		_, wrapError := wrapper.WrapValidateError(op, vErr, entity.Meta{"request": req})
		return wrapError
	}

	if req.Amount.Equal(entity.Amount{}) || req.Amount.LessThan(entity.Amount{}) {
		return richerror.New(op, errors.New("locked price must be greater than zero")).
			WithKind(richerror.KindValidation).
			WithCode(richerror.CodeInvalidInput).
			WithMessage("locked price must be greater than zero")
	}

	if req.PriceLockedAt <= 0 || req.PriceLockExpiresAt <= 0 {
		return richerror.New(op, errors.New("price lock metadata is invalid")).
			WithKind(richerror.KindValidation).
			WithCode(richerror.CodeInvalidInput).
			WithMessage("price lock metadata is invalid")
	}

	if req.PriceLockExpiresAt <= time.Now().Unix() {
		return richerror.New(op, errors.New("price lock has expired")).
			WithKind(richerror.KindValidation).
			WithCode(richerror.CodeInvalidInput).
			WithMessage("price lock has expired")
	}

	if req.PriceLockPaymentMethod != method {
		return richerror.New(op, errors.New("price lock payment method mismatch")).
			WithKind(richerror.KindValidation).
			WithCode(richerror.CodeInvalidInput).
			WithMessage("price lock payment method mismatch")
	}

	return nil
}
