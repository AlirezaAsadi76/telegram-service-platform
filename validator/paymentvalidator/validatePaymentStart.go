package paymentvalidator

import (
	"errors"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func (v *Validator) ValidateStartPayment(req paymentparams.StartPaymentHandlerRequest) (map[string]string, error) {
	const Op = "validator.validateStartPayment"

	vErr := validation.ValidateStruct(&req,

		validation.Field(&req.OrderID, validation.Required, validation.Min(uint64(1))),

		validation.Field(&req.Method, validation.Required,
			validation.In(paymententity.PaymentMethodZarinpal, paymententity.PaymentMethodCrypto)),

		validation.Field(&req.Amount, validation.Required),

		validation.Field(&req.Currency,
			validation.Required,
			validation.In(entity.CurrencyTOMAN, entity.CurrencyUSD, entity.CurrencyTON, entity.CurrencyUSDT),
		),

		validation.Field(&req.IdempotencyKey, validation.Required),

		validation.Field(&req.CallbackURL,
			validation.Required.Error("callback_url is required"), is.URL),

		validation.Field(&req.Description,
			validation.Length(0, 255).Error("description must be at most 255 characters"),
		),
	)

	if vErr != nil {

		fieldErrors := make(map[string]string)
		var errV validation.Errors
		ok := errors.As(vErr, &errV)
		if ok {
			for key, val := range errV {
				fieldErrors[key] = val.Error()
			}
		}
		return fieldErrors, richerror.New(Op, vErr).
			WithMessage(msgerror.InvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"request": req})
	}

	return nil, nil

}
