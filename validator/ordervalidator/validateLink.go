package ordervalidator

import (
	"regexp"
	"telegram-service-platform/entity"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/pkg/wrapper"

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
		_, wrapError := wrapper.WrapValidateError(op, vErr, entity.Meta{"request": req})
		return wrapError
	}
	return nil
}
