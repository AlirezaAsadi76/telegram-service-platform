package postgrescheckout

import (
	"context"
	"errors"
	"fmt"

	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/repository/postgresorder"
	"telegram-service-platform/repository/postgreswallet"

	"github.com/jackc/pgx/v5"
)

const IdempotencyRefund = "refund:order:%d"

func (d *DB) ExecuteOrderRefund(ctx context.Context, req checkoutparams.RefundOrderRequest) (*checkoutparams.RefundOrderResponse, error) {
	const Op = "postgrescheckout.ExecuteOrderRefund"

	var result checkoutparams.RefundOrderResponse

	err := d.transactionProvider.Execute(ctx, func(tx pgx.Tx) error {
		walletRepo := postgreswallet.NewWithExecutor(tx)
		orderRepo := postgresorder.NewWithExecutor(tx)

		order, err := orderRepo.GetByIDForUpdate(ctx, req.OrderID)
		if err != nil {
			return err
		}

		if order.Status == orderentity.OrderStatusFailed {
			existingRefund, refundErr := walletRepo.GetTransactionByIdempotencyKey(ctx,
				fmt.Sprintf(IdempotencyRefund, order.ID),
			)

			if refundErr == nil && existingRefund != nil {
				result = checkoutparams.RefundOrderResponse{
					OrderID:         order.ID,
					UserID:          existingRefund.UserID,
					WalletTxID:      existingRefund.ID,
					RefundAmount:    existingRefund.Amount,
					AlreadyRefunded: true,
				}

				return nil
			}

			return richerror.New(Op,
				fmt.Errorf("order %d is failed but refund transaction was not found", order.ID)).
				WithKind(richerror.KindConflict).
				WithMessage(msgerror.OrderUpdateFailed)
		}

		if order.Status != orderentity.OrderStatusProcessing {
			return richerror.New(Op, fmt.Errorf("order %d cannot be refunded from status %s", order.ID, order.Status)).
				WithKind(richerror.KindConflict).
				WithCode(richerror.CodeOrderInvalidState).
				WithMessage(msgerror.OrderUpdateFailed)
		}

		wallet, gErr := walletRepo.GetForUpdate(ctx, order.UserID)
		if gErr != nil {
			return gErr
		}

		idempotencyKey := fmt.Sprintf(IdempotencyRefund, order.ID)

		existingTx, gtErr := walletRepo.GetTransactionByIdempotencyKey(ctx, idempotencyKey)

		if gtErr == nil && existingTx != nil {
			return richerror.New(Op, fmt.Errorf("refund transaction exists for order %d but order is still processing", order.ID)).
				WithKind(richerror.KindConflict).
				WithMessage(msgerror.OrderUpdateFailed)
		}

		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		refundTx := &walletentity.WalletTransaction{
			WalletID:       wallet.ID,
			UserID:         order.UserID,
			Type:           walletentity.WalletTransactionTypeRefund,
			Amount:         order.Amount,
			Status:         walletentity.WalletTransactionStatusPending,
			ReferenceID:    fmt.Sprintf(IdempotencyRefund, order.ID),
			IdempotencyKey: idempotencyKey,
		}

		if err := walletRepo.CreateTransaction(ctx, refundTx); err != nil {
			return err
		}

		newBalance := wallet.Balance.Add(order.Amount)

		if err := walletRepo.UpdateBalanceAtomic(
			ctx,
			wallet.ID,
			newBalance,
			wallet.Version+1,
		); err != nil {
			return err
		}

		if err := walletRepo.UpdateTransactionStatus(ctx, refundTx.ID,
			walletentity.WalletTransactionStatusComplete,
		); err != nil {
			return err
		}

		if err := orderRepo.UpdateStatus(
			ctx,
			order.ID,
			orderentity.OrderStatusFailed,
			"",
			nil,
		); err != nil {
			return err
		}

		result = checkoutparams.RefundOrderResponse{
			UserID:       order.UserID,
			OrderID:      order.ID,
			WalletTxID:   refundTx.ID,
			RefundAmount: order.Amount,
		}

		return nil
	})

	if err != nil {
		if richerror.IsKind(err, richerror.KindConflict) {
			return nil, err
		}

		return nil, richerror.New(Op, err).
			WithKind(richerror.KindQueryFailure).
			WithMessage(msgerror.OrderUpdateFailed)
	}

	return &result, nil
}
