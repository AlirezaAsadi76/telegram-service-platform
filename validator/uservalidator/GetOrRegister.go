package uservalidator

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/params/userparams"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/pkg/wrapper"

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
		_, wrapError := wrapper.WrapValidateError(op, vErr,
			entity.Meta{
				"telegram_id": req.TelegramID,
			})
		return wrapError
	}

	return nil
}
