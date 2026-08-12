package checkoutservice

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/paymentparams"
	"time"
)

func (s *Service) ProcessWalletRecharge(ctx context.Context, req checkoutparams.ProcessWalletRechargeRequest) (
	*checkoutparams.ProcessWalletRechargeResponse, error) {
	// Generate idempotency key
	keyData := fmt.Sprintf("recharge:%d:%d:%d", req.UserID, req.Amount, time.Now().Unix()/60) // per-minute granularity
	hash := sha256.Sum256([]byte(keyData))
	idempotencyKey := hex.EncodeToString(hash[:])

	// Create payment
	createResp, err := s.paymentSvc.Create(ctx, paymentparams.CreateRequest{
		UserID:         req.UserID,
		Method:         req.Method,
		Amount:         req.Amount,
		Currency:       req.Currency,
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return &checkoutparams.ProcessWalletRechargeResponse{
		PaymentID:  createResp.PaymentID,
		PaymentURL: createResp.PaymentURL,
	}, nil
}
