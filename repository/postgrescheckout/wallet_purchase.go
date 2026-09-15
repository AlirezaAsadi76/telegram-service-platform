package postgrescheckout

import (
	"context"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/params/walletparam"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/repository/postgresorder"
	"telegram-service-platform/repository/postgreswallet"

	"github.com/jackc/pgx/v5"
)

func (d *DB) ExecuteWalletPurchase(
	ctx context.Context,
	req walletparam.WalletPurchaseRequest,
) (*walletparam.WalletPurchaseResult, error) {
	const Op = "postgrescheckout.ExecuteWalletPurchase"

	var result walletparam.WalletPurchaseResult

	err := d.transactionProvider.Execute(ctx, func(tx pgx.Tx) error {
		walletRepo := postgreswallet.NewWithExecutor(tx)
		orderRepo := postgresorder.NewWithExecutor(tx)

		wallet, err := walletRepo.GetForUpdate(ctx, req.UserID)
		if err != nil {
			return err
		}

		if !wallet.HasSufficient(req.Amount) {
			return richerror.New(Op, nil).
				WithKind(richerror.KindValidation).
				WithMessage(msgerror.InsufficientBalance)
		}

		order := &orderentity.Order{
			UserID:      req.UserID,
			ProductType: req.ProductType,
			ProductID:   req.ProductID,
			Quantity:    req.Quantity,
			TargetLink:  req.TargetLink,
			Amount:      req.Amount,
			Currency:    req.Currency,
			Status:      orderentity.OrderStatusPaid,
		}

		if err := orderRepo.Create(ctx, order); err != nil {
			return err
		}

		walletTx := &walletentity.WalletTransaction{
			WalletID:       wallet.ID,
			UserID:         req.UserID,
			Type:           walletentity.WalletTransactionTypeWithdraw,
			Amount:         req.Amount,
			Status:         walletentity.WalletTransactionStatusPending,
			ReferenceID:    req.ReferenceID,
			IdempotencyKey: req.IdempotencyKey,
		}

		if err := walletRepo.CreateTransaction(ctx, walletTx); err != nil {
			return err
		}

		newBalance := wallet.Balance.Sub(req.Amount)

		if err := walletRepo.UpdateBalanceAtomic(
			ctx,
			wallet.ID,
			newBalance,
			wallet.Version+1,
		); err != nil {
			return err
		}

		if err := walletRepo.UpdateTransactionStatus(
			ctx,
			walletTx.ID,
			walletentity.WalletTransactionStatusComplete,
		); err != nil {
			return err
		}

		result = walletparam.WalletPurchaseResult{
			OrderID:     order.ID,
			WalletTxID:  walletTx.ID,
			NewBalance:  newBalance,
			OrderStatus: orderentity.OrderStatusPaid,
		}

		return nil
	})

	if err != nil {
		if richerror.IsKind(err, richerror.KindValidation) {
			return nil, err
		}

		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.OrderUpdateFailed)
	}

	return &result, nil
}
