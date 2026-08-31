package uservalidator

import (
	"context"
	"errors"
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/params"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidationUserExistence(ctx context.Context, telegramID entity.TelegramId) (uint64, error) {
	const op richerror.Op = "uservalidator.ValidationUserExistence"

	vErr := validation.Validate(telegramID,
		validation.Required,
		validation.Min(int64(1)),
		validation.By(v.ensureUserExists(ctx)),
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
			return 0, richerror.New(op, vErr).
				WithKind(richerror.KindValidation).
				WithMessage(firstErrorMsg).
				WithMeta(map[string]interface{}{"telegram_id": telegramID})
		}
	}

	user, _ := v.userService.FindUserByTelegramID(ctx, params.FindUserByTelegramIDRequest{
		TelegramID: telegramID,
	})

	return user.UserInfo.Id, nil
}

func (v Validator) ensureUserExists(ctx context.Context) validation.RuleFunc {
	return func(value interface{}) error {
		telegramID := value.(int64)

		user, err := v.userService.FindUserByTelegramID(ctx, params.FindUserByTelegramIDRequest{
			TelegramID: entity.TelegramId(telegramID),
		})

		if err != nil || !user.Found {
			return fmt.Errorf(msgerror.UserNotFound)
		}

		return nil
	}
}
