package ordervalidator

import (
	"errors"
	"regexp"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
)

var urlPattern = regexp.MustCompile(`^https?://`)

func (v *Validator) ValidateLink(req orderparams.SubmitLinkRequest) error {
	const op richerror.Op = "ordervalidator.ValidateLink"

	vErr := validation.ValidateStruct(&req,
		validation.Field(&req.Link,
			validation.Required.Error("لینک نمی‌تواند خالی باشد"),
			validation.Match(urlPattern).Error("لینک نامعتبر است. لطفاً لینک را با http:// یا https:// شروع کنید"),
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
