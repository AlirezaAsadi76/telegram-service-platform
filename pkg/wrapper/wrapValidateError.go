package wrapper

import (
	"errors"
	"telegram-service-platform/entity"
	"telegram-service-platform/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
)

func WrapValidateError(operation richerror.Op, validateError error, meta entity.Meta) (entity.Meta, error) {
	var errV validation.Errors
	if ok := errors.As(validateError, &errV); ok {
		fieldErrors := make(entity.Meta)
		firstErrorMsg := ""
		for key, val := range errV {
			if firstErrorMsg == "" {
				firstErrorMsg = val.Error()
			}
			fieldErrors[key] = val.Error()
		}

		return fieldErrors, richerror.New(operation, validateError).
			WithKind(richerror.KindValidation).
			WithCode(richerror.CodeInvalidInput).
			WithMessage(firstErrorMsg).
			WithMeta(meta)
	}

	return nil, nil
}
