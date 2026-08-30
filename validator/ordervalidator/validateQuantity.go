package ordervalidator

import (
	"errors"
	"fmt"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v *Validator) ValidateQuantity(req orderparams.SubmitQuantityRequest) error {
	const op richerror.Op = "ordervalidator.ValidateQuantity"

	vErr := validation.ValidateStruct(&req,
		validation.Field(&req.Quantity,
			validation.Required.Error("تعداد نمی‌تواند خالی باشد"),
			validation.Min(req.Min).Error(fmt.Sprintf("تعداد سفارش باید حداقل %d باشد", req.Min)),
			validation.Max(req.Max).Error(fmt.Sprintf("تعداد سفارش نباید بیشتر از %d باشد", req.Max)),
		),
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
				WithMeta(map[string]interface{}{
					"request":      req,
					"field_errors": errV,
				})
		}
	}
	return nil
}
