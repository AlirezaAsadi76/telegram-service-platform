package checkoutservice

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/orderparams"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/params/walletparam"
	"time"
)

func (s *Service) ProcessProductPurchase(ctx context.Context, req checkoutparams.ProcessProductPurchaseRequest) error {
	// 1. Create order
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
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Generate idempotency key
	keyData := fmt.Sprintf("purchase:%d:%d:%d", req.UserID, orderResp.OrderID, time.Now().Unix())
	hash := sha256.Sum256([]byte(keyData))
	idempotencyKey := hex.EncodeToString(hash[:])

	if req.Method == paymententity.PaymentMethodManual {
		// Pay from wallet
		debitReq := walletparam.DebitRequest{
			UserID:         req.UserID,
			Amount:         req.Amount,
			ReferenceID:    fmt.Sprintf("order:%d", orderResp.OrderID),
			IdempotencyKey: idempotencyKey,
		}
		_, err := s.walletSvc.Debit(ctx, debitReq)
		if err != nil {
			// Cancel order
			s.orderSvc.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
				OrderID: orderResp.OrderID,
				Status:  orderentity.OrderStatusCanceled,
			})
			return fmt.Errorf("insufficient wallet balance: %w", err)
		}

		// Mark as paid and fulfill
		s.orderSvc.UpdateStatus(ctx, orderparams.UpdateStatusRequest{
			OrderID: orderResp.OrderID,
			Status:  orderentity.OrderStatusPaid,
		})

		// Fulfill via SMM provider (async - push to queue)
		// In real implementation, this would be a background job
		order, _ := s.orderSvc.GetById(ctx, orderResp.OrderID)
		go s.fulfillOrderAsync(order)

	} else {
		// External payment
		paymentResp, err := s.paymentSvc.Create(ctx, paymentparams.CreateRequest{
			OrderID:        &orderResp.OrderID,
			UserID:         req.UserID,
			Method:         req.Method,
			Amount:         req.Amount,
			Currency:       req.Currency,
			IdempotencyKey: idempotencyKey,
		})
		if err != nil {
			return fmt.Errorf("failed to create payment: %w", err)
		}

		// Send payment link to user
		s.messenger.SendToUser(ctx, int64(req.UserID), fmt.Sprintf("Please complete payment: %s", paymentResp.PaymentURL))
	}

	return nil
}
