package checkoutservice

import (
	"context"
	"fmt"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/params/walletparam"
	"telegram-service-platform/pkg/hashing"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/pkg/ts"
)

func (s *Service) ProcessManualWalletRecharge(ctx context.Context, req checkoutparams.ManualRechargeRequest) error {
	const Op = "checkoutservice.ProcessManualWalletRecharge"

	idempotencyKey := hashing.EncodeStringToSHA256(
		fmt.Sprintf("%s:%d:%d:%d:%d", s.config.PrefixManualIdempotencyKey, req.AdminID, req.UserID, req.Amount, ts.Now()))

	_, err := s.walletSvc.Credit(ctx, walletparam.CreditRequest{
		UserID:         req.UserID,
		Amount:         req.Amount,
		ReferenceID:    fmt.Sprintf("manual_admin:%d", req.AdminID),
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		return richerror.New(Op, err)
	}

	_ = s.messenger.SendToUser(ctx, int64(req.UserID),
		fmt.Sprintf("Your wallet charged: %d by admin", req.Amount))
	_ = s.messenger.SendToAdminChannel(ctx,
		fmt.Sprintf("Admin %d recharged user %d: %d", req.AdminID, req.UserID, req.Amount))

	return nil
}
