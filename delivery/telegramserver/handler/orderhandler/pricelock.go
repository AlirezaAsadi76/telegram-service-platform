package orderhandler

import (
	"context"
	"errors"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/entity/productentity"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/orderparams"
	"time"
)

var ErrPriceLockMethodMismatch = errors.New(
	"price lock payment method mismatch",
)

func (h *Handler) ensurePriceLock(
	ctx context.Context,
	telegramID int64,
	userID uint64,
	state *orderentity.OrderFlowState,
	method entity.PriceLockPaymentMethod,
) error {
	now := time.Now().Unix()

	if state.PriceLockExpiresAt > now {
		if state.PriceLockPaymentMethod != string(method) {
			return ErrPriceLockMethodMismatch
		}

		return nil
	}
	req := checkoutparams.LockProductPriceRequest{
		UserID:        userID,
		ProductType:   productentity.ProductTypeSMM,
		ProductID:     state.ServiceID,
		Quantity:      state.Quantity,
		Currency:      state.Currency,
		PaymentMethod: method,
	}
	if vErr := h.validator.ValidateLockProductPrice(req); vErr != nil {
		return vErr
	}
	response, err := h.checkoutService.LockProductPrice(ctx, req)
	if err != nil {
		return err
	}

	state.Price = response.Amount
	state.Currency = response.Currency
	state.PriceLockedAt = response.LockedAt
	state.PriceLockExpiresAt = response.ExpiresAt
	state.PriceLockPaymentMethod = string(response.PaymentMethod)

	// Explicit phase transition:
	// selection TTL → payment/price-lock TTL.
	state.ExpiresAt = response.ExpiresAt

	return h.orderFlowService.SaveOrderFlow(
		ctx,
		orderparams.SaveOrderFlowRequest{
			TelegramID: entity.TelegramId(telegramID),
			State:      *state,
		},
	)
}

func paymentMethodToPriceLockMethod(
	method paymententity.PaymentMethod,
) entity.PriceLockPaymentMethod {
	return entity.PriceLockPaymentMethod(method)
}
