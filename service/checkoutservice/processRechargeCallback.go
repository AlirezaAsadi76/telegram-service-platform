package checkoutservice

import (
	"context"
	"fmt"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/params/walletparam"
)

func (s *Service) ProcessRechargeCallback(ctx context.Context, paymentID uint64, externalID string, callbackData map[string]any) error {
	// 1. Verify payment
	verifyResp, err := s.paymentSvc.Verify(ctx, paymentparams.VerifyRequest{
		PaymentID:    paymentID,
		ExternalID:   externalID,
		CallbackData: callbackData,
	})
	if err != nil {
		return fmt.Errorf("payment verification failed: %w", err)
	}

	if verifyResp.Status != paymententity.PaymentStatusSuccess {
		// Payment failed or pending - notify user
		s.messenger.SendToUser(ctx, int64(paymentID), "Payment failed. Please try again.")
		return nil
	}

	// 2. Get payment details to find amount
	payment, _ := s.paymentSvc.GetById(ctx, paymentID)

	// 3. Credit wallet (idempotent)
	creditReq := walletparam.CreditRequest{
		UserID:         payment.UserID,
		Amount:         payment.Amount,
		ReferenceID:    fmt.Sprintf("payment:%d", payment.ID),
		IdempotencyKey: payment.IdempotencyKey,
	}
	_, err = s.walletSvc.Credit(ctx, creditReq)
	if err != nil {
		return fmt.Errorf("wallet credit failed: %w", err)
	}

	// 4. Notify
	s.messenger.SendToUser(ctx, int64(payment.UserID), fmt.Sprintf("Wallet credited: %d %s", payment.Amount, payment.Currency))
	s.messenger.SendToAdminChannel(ctx, fmt.Sprintf("User %d recharged wallet: %d %s", payment.UserID, payment.Amount, payment.Currency))

	return nil
}
