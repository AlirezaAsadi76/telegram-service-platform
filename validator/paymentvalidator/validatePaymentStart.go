package paymentvalidator

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/wrapper"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func (v *Validator) ValidateStartPayment(req paymentparams.StartPaymentHandlerRequest) (entity.Meta, error) {
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

		fieldErrors, wrapError := wrapper.WrapValidateError(Op, vErr, entity.Meta{"request": req})
		return fieldErrors, wrapError

	}

	return nil, nil

}
