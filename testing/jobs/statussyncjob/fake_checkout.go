package statussynctesting

import (
	"context"
	"telegram-service-platform/params/checkoutparams"

	"telegram-service-platform/params/walletparam"
)

type fakeCheckoutTransactionRepository struct {
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

	f.calls = append(f.calls, req)

	if f.err != nil {
		return nil, f.err
	}

	return &checkoutparams.RefundOrderResponse{
		OrderID: req.OrderID,
		UserID:  100,
	}, nil
}
