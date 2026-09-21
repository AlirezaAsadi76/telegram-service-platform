package checkouttesting

import (
	"context"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/walletparam"
)

type fakeTransactionRepository struct {
	calls []checkoutparams.RefundOrderRequest

	response *checkoutparams.RefundOrderResponse
	err      error
}

func (f *fakeTransactionRepository) ExecuteWalletPurchase(
	_ context.Context,
	_ walletparam.WalletPurchaseRequest,
) (*walletparam.WalletPurchaseResult, error) {
	return nil, nil
}

func (f *fakeTransactionRepository) ExecuteOrderRefund(
	_ context.Context,
	req checkoutparams.RefundOrderRequest,
) (*checkoutparams.RefundOrderResponse, error) {
	f.calls = append(f.calls, req)

	if f.err != nil {
		return nil, f.err
	}

	return f.response, nil
}
