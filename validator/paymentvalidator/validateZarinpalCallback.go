package paymentvalidator

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/wrapper"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v *Validator) ValidateZarinpalCallback(req paymentparams.ZarinpalCallback) (entity.Meta, error) {
	const Op = "validate.ZarinpalCallback"

	vErr := validation.ValidateStruct(&req,
		validation.Field(&req.Authority, validation.Required, validation.NilOrNotEmpty),
		validation.Field(&req.Status, validation.Required, validation.NilOrNotEmpty, validation.In("OK", "NOK")))

	if vErr != nil {

		fieldErrors, wrapError := wrapper.WrapValidateError(Op, vErr, entity.Meta{"request": req})
		return fieldErrors, wrapError

	}

	return nil, nil

}
