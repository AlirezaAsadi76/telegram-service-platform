package orderfulfillmenttesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/walletparam"
)

type fakeCheckoutTransactionRepository struct {
	order *orderentity.Order

	calls []checkoutparams.RefundOrderRequest

	err error
}

func (f *fakeCheckoutTransactionRepository) ExecuteWalletPurchase(
	_ context.Context,
	_ walletparam.WalletPurchaseRequest,
) (*walletparam.WalletPurchaseResult, error) {
	return nil, nil
}

func (f *fakeCheckoutTransactionRepository) ExecuteOrderRefund(
	_ context.Context,
	req checkoutparams.RefundOrderRequest,
) (*checkoutparams.RefundOrderResponse, error) {
	f.calls = append(
		f.calls,
		req,
	)

	if f.err != nil {
		return nil, f.err
	}

	if f.order == nil {
		return nil, errors.New("order not found")
	}

	if f.order.Status != orderentity.OrderStatusProcessing {
		return nil, errors.New("order is not processing")
	}

	f.order.Status = orderentity.OrderStatusFailed

	return &checkoutparams.RefundOrderResponse{
		OrderID: req.OrderID,
		UserID:  f.order.UserID,
	}, nil
}
