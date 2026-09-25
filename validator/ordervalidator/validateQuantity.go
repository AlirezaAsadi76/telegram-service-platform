package ordervalidator

import (
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/pkg/wrapper"

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

		_, wrapError := wrapper.WrapValidateError(op, vErr, entity.Meta{"request": req})
		return wrapError

	}
	return nil
}
