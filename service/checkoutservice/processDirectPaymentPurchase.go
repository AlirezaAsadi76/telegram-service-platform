package checkoutservice

import (
	"context"
	"fmt"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/hashing"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/pkg/ts"
)

func (s *Service) ProcessDirectPaymentPurchase(ctx context.Context, req checkoutparams.DirectPaymentPurchase) (*checkoutparams.PaymentURLResponse, error) {
	const Op = "checkoutservice.ProcessDirectPaymentPurchase"

	// 1. Create Order (PENDING)
	orderResp, err := s.orderSvc.Create(ctx, orderparams.CreateRequest{
		UserID:      req.UserID,
		ProductType: req.ProductType,
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		TargetLink:  req.TargetLink,
		Amount:      req.Amount,
		Currency:    req.Currency,
	})
	if err != nil {
		return nil, richerror.New(Op, err)
	}

	// 2. Generate idempotency key
	idempotencyKey := hashing.EncodeStringToSHA256(
		fmt.Sprintf("%s:%d:%d:%d", s.config.PrefixDirectIdempotencyKey, req.UserID, orderResp.OrderID, ts.Now()),
	)

	// 3. Create Payment
	paymentResp, cpErr := s.paymentSvc.Create(ctx, paymentparams.CreateRequest{
		OrderID:        orderResp.OrderID,
		UserID:         req.UserID,
		Method:         req.Method,
		Amount:         req.Amount,
		Currency:       req.Currency,
		IdempotencyKey: idempotencyKey,
	})
	if cpErr != nil {
		return nil, richerror.New(Op, cpErr)
	}

	return &checkoutparams.PaymentURLResponse{
		OrderID:    orderResp.OrderID,
		PaymentID:  paymentResp.PaymentID,
		PaymentURL: paymentResp.PaymentURL,
	}, nil
}
