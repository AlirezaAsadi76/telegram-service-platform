package checkoutValidator

import (
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/pkg/wrapper"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v *Validator) ValidateLockSMMPrice(req checkoutparams.LockSMMPriceRequest) error {
	const op = "checkout.ValidateLockSMMPrice"

	vErr := validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required, validation.Min(1)),
		validation.Field(&req.ProductID, validation.Required, validation.Min(1)),
		validation.Field(&req.Quantity, validation.Required, validation.Min(1)),
		validation.Field(&req.ProductType, validation.Required, validation.In(
			productentity.ProductTypeSMM, productentity.ProductTypeAds, productentity.ProductTypePremium, productentity.ProductTypeStars)),
	)

	if vErr != nil {
		_, wrapError := wrapper.WrapValidateError(op, vErr, entity.Meta{"request": req})
		return wrapError
	}

	return nil
}
