package uservalidator

import (
	"telegram-service-platform/params"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (v Validator) GetOrRegister(
	req params.GetOrRegisterRequest,
) error {

	const op richerror.Op = "uservalidator.GetOrRegister"

	if req.TelegramID <= 0 {

		return richerror.New(op, nil).
			WithKind(richerror.KindInvalid).
			WithMessage(msgerror.TelegramIdInvalid)
	}
	return nil
}
