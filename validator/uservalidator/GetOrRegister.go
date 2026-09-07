package uservalidator

import (
	"errors"
	"telegram-service-platform/params/userparams"
	"telegram-service-platform/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) GetOrRegister(
	req userparams.GetOrRegisterRequest,
) error {

	const op richerror.Op = "uservalidator.GetOrRegister"

	vErr := validation.ValidateStruct(&req,
		validation.Field(&req.TelegramID, validation.Required, validation.NilOrNotEmpty, validation.Min(0)),
		validation.Field(&req.Role, validation.Required),
	)

	if vErr != nil {
		var errV validation.Errors
		if ok := errors.As(vErr, &errV); ok {
			var firstErrorMsg string
			for _, err := range errV {
				if err != nil {
					firstErrorMsg = err.Error()
					break
				}
			}
			return richerror.New(op, vErr).
				WithKind(richerror.KindValidation).
				WithMessage(firstErrorMsg).
				WithMeta(map[string]interface{}{"telegram_id": req.TelegramID})
		}
	}

	return nil
}
