package paymentvalidator

import (
	"errors"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v *Validator) ValidateZarinpalCallback(req paymentparams.ZarinpalCallback) (map[string]string, error) {
	const Op = "validate.ZarinpalCallback"

	vErr := validation.ValidateStruct(&req,
		validation.Field(&req.Authority, validation.Required, validation.NilOrNotEmpty),
		validation.Field(&req.Status, validation.Required, validation.NilOrNotEmpty, validation.In("OK", "NOK")))

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
