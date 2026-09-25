package uservalidator

import (
	"context"
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/params/userparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/pkg/wrapper"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidationUserExistence(ctx context.Context, telegramID entity.TelegramId) (uint64, error) {
	const op richerror.Op = "uservalidator.ValidationUserExistence"

	vErr := validation.Validate(telegramID.Int64(),
		validation.Required,
		validation.Min(int64(1)),
		validation.By(v.ensureUserExists(ctx)),
	)

	if vErr != nil {
		_, richError := wrapper.WrapValidateError(op, vErr, entity.Meta{"telegram_id": telegramID})
		return 0, richError
	}

	user, _ := v.userService.FindUserByTelegramID(ctx, userparams.FindUserByTelegramIDRequest{
		TelegramID: telegramID,
	})

	return user.UserInfo.Id, nil
}

func (v Validator) ensureUserExists(ctx context.Context) validation.RuleFunc {
	return func(value interface{}) error {
		telegramID := value.(int64)

		user, err := v.userService.FindUserByTelegramID(ctx, userparams.FindUserByTelegramIDRequest{
			TelegramID: entity.TelegramId(telegramID),
		})

		if err != nil || !user.Found {
			return fmt.Errorf(msgerror.UserNotFound)
		}

		return nil
	}
}
