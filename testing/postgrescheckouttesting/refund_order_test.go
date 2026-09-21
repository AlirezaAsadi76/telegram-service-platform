package postgrescheckouttesting

import (
	"context"
	"fmt"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/walletentity"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/pkg/richerror"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgrescheckout"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

func newTestCheckoutRepository(
	t *testing.T,
) (*postgrescheckout.DB, *pgxpool.Pool) {
	t.Helper()

	pool := newTestPool(t)

	transactionProvider := postgres.NewTransactionProvider(pool)

	checkoutRepo := postgrescheckout.New(
		transactionProvider,
	)

	return checkoutRepo, pool
}

func TestExecuteOrderRefund_Success(t *testing.T) {

	checkoutRepo, pool := newTestCheckoutRepository(t)

	refundAmount := decimal.NewFromInt(30)

	fixture := createRefundFixture(
		t,
		pool,
		orderentity.OrderStatusProcessing,
		refundAmount,
	)

	beforeBalance := getWalletBalance(
		t,
		pool,
		fixture.WalletID,
	)

	result, err := checkoutRepo.ExecuteOrderRefund(
		context.Background(),
		checkoutparams.RefundOrderRequest{
			OrderID: fixture.OrderID,
			Reason:  "provider_failed",
		},
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if result == nil {
		t.Fatal("expected refund result")
	}

	if result.AlreadyRefunded {
		t.Fatal("expected first refund not to be marked AlreadyRefunded")
	}

	if result.OrderID != fixture.OrderID {
		t.Fatalf(
			"expected order ID %d, got %d",
			fixture.OrderID,
			result.OrderID,
		)
	}

	if result.UserID != fixture.UserID {
		t.Fatalf(
			"expected user ID %d, got %d",
			fixture.UserID,
			result.UserID,
		)
	}

	if !result.RefundAmount.Equal(entity.Amount(refundAmount)) {
		t.Fatalf(
			"expected refund amount %s, got %s",
			refundAmount,
			result.RefundAmount,
		)
	}

	afterBalance := getWalletBalance(
		t,
		pool,
		fixture.WalletID,
	)

	expectedBalance := beforeBalance.Add(refundAmount)

	if !afterBalance.Equal(expectedBalance) {
		t.Fatalf(
			"expected wallet balance %s, got %s",
			expectedBalance,
			afterBalance,
		)
	}

	status, providerID, externalOrderID := getOrderState(
		t,
		pool,
		fixture.OrderID,
	)

	if status != orderentity.OrderStatusFailed {
		t.Fatalf(
			"expected order status FAILED, got %s",
			status,
		)
	}

	if providerID == nil {
		t.Fatal("expected provider ID to be preserved")
	}

	if *providerID != fixture.ProviderID {
		t.Fatalf(
			"expected provider ID %d, got %d",
			fixture.ProviderID,
			*providerID,
		)
	}

	if externalOrderID == "" {
		t.Fatal("expected external order ID to be preserved")
	}

	resultRefund := getRefundTransaction(t, pool, fixture.OrderID)

	if resultRefund.txID == 0 {
		t.Fatal("expected refund transaction ID")
	}

	if resultRefund.txType != string(walletentity.WalletTransactionTypeRefund) {
		t.Fatalf(
			"expected transaction type %s, got %s",
			walletentity.WalletTransactionTypeRefund,
			resultRefund.txType,
		)
	}

	if !resultRefund.amount.Equal(entity.Amount(refundAmount)) {
		t.Fatalf(
			"expected transaction amount %s, got %s",
			refundAmount,
			resultRefund.amount,
		)
	}

	if resultRefund.status != string(walletentity.WalletTransactionStatusComplete) {
		t.Fatalf(
			"expected transaction status %s, got %s",
			walletentity.WalletTransactionStatusComplete,
			resultRefund.status,
		)
	}

	expectedKey := fmt.Sprintf(
		postgrescheckout.IdempotencyRefund,
		fixture.OrderID,
	)

	if resultRefund.idempotency_key != expectedKey {
		t.Fatalf(
			"expected idempotency key %s, got %s",
			expectedKey,
			resultRefund.idempotency_key,
		)
	}

	if resultRefund.referenceID != expectedKey {
		t.Fatalf(
			"expected reference ID %s, got %s",
			expectedKey,
			resultRefund.referenceID,
		)
	}
}

func TestExecuteOrderRefund_Idempotent(t *testing.T) {

	checkoutRepo, pool := newTestCheckoutRepository(t)

	refundAmount := decimal.NewFromInt(30)

	fixture := createRefundFixture(
		t,
		pool,
		orderentity.OrderStatusProcessing,
		refundAmount,
	)

	req := checkoutparams.RefundOrderRequest{
		OrderID: fixture.OrderID,
		Reason:  "provider_failed",
	}

	firstResult, firstErr := checkoutRepo.ExecuteOrderRefund(
		context.Background(),
		req,
	)

	if firstErr != nil {
		t.Fatalf(
			"unexpected first refund error: %v",
			firstErr,
		)
	}

	balanceAfterFirst := getWalletBalance(
		t,
		pool,
		fixture.WalletID,
	)

	txCountAfterFirst := getRefundTransactionCount(
		t,
		pool,
		fixture.OrderID,
	)

	secondResult, secondErr := checkoutRepo.ExecuteOrderRefund(
		context.Background(),
		req,
	)

	if secondErr != nil {
		t.Fatalf(
			"unexpected second refund error: %v",
			secondErr,
		)
	}

	if firstResult == nil {
		t.Fatal("expected first result")
	}

	if secondResult == nil {
		t.Fatal("expected second result")
	}

	if firstResult.AlreadyRefunded {
		t.Fatal("expected first refund to be fresh")
	}

	if !secondResult.AlreadyRefunded {
		t.Fatal("expected second refund to be AlreadyRefunded")
	}

	balanceAfterSecond := getWalletBalance(
		t,
		pool,
		fixture.WalletID,
	)

	if !balanceAfterSecond.Equal(balanceAfterFirst) {
		t.Fatalf(
			"expected balance to remain %s, got %s",
			balanceAfterFirst,
			balanceAfterSecond,
		)
	}

	txCountAfterSecond := getRefundTransactionCount(
		t,
		pool,
		fixture.OrderID,
	)

	if txCountAfterFirst != 1 {
		t.Fatalf(
			"expected one refund transaction after first call, got %d",
			txCountAfterFirst,
		)
	}

	if txCountAfterSecond != 1 {
		t.Fatalf(
			"expected one refund transaction after second call, got %d",
			txCountAfterSecond,
		)
	}

	status, providerID, externalOrderID := getOrderState(
		t,
		pool,
		fixture.OrderID,
	)

	if status != orderentity.OrderStatusFailed {
		t.Fatalf(
			"expected order status FAILED, got %s",
			status,
		)
	}

	if providerID == nil {
		t.Fatal("expected provider ID to remain preserved")
	}

	if *providerID != fixture.ProviderID {
		t.Fatalf(
			"expected provider ID %d, got %d",
			fixture.ProviderID,
			*providerID,
		)
	}

	if externalOrderID == "" {
		t.Fatal("expected external order ID to remain preserved")
	}
}

func TestExecuteOrderRefund_InvalidState(t *testing.T) {

	checkoutRepo, pool := newTestCheckoutRepository(t)

	fixture := createRefundFixture(
		t,
		pool,
		orderentity.OrderStatusPaid,
		decimal.NewFromInt(30),
	)

	beforeBalance := getWalletBalance(
		t,
		pool,
		fixture.WalletID,
	)

	_, err := checkoutRepo.ExecuteOrderRefund(
		context.Background(),
		checkoutparams.RefundOrderRequest{
			OrderID: fixture.OrderID,
			Reason:  "invalid_state",
		},
	)

	if err == nil {
		t.Fatal("expected invalid state error")
	}

	if !richerror.IsKind(
		err,
		richerror.KindConflict,
	) {
		t.Fatalf(
			"expected conflict error, got %v",
			err,
		)
	}

	if !richerror.IsCode(
		err,
		richerror.CodeOrderInvalidState,
	) {
		t.Fatalf(
			"expected error code %s, got %v",
			richerror.CodeOrderInvalidState,
			err,
		)
	}

	afterBalance := getWalletBalance(
		t,
		pool,
		fixture.WalletID,
	)

	if !afterBalance.Equal(beforeBalance) {
		t.Fatalf(
			"expected wallet balance to remain %s, got %s",
			beforeBalance,
			afterBalance,
		)
	}

	txCount := getRefundTransactionCount(
		t,
		pool,
		fixture.OrderID,
	)

	if txCount != 0 {
		t.Fatalf(
			"expected no refund transaction, got %d",
			txCount,
		)
	}

	status, providerID, externalOrderID := getOrderState(
		t,
		pool,
		fixture.OrderID,
	)

	if status != orderentity.OrderStatusPaid {
		t.Fatalf(
			"expected order status PAID, got %s",
			status,
		)
	}

	if providerID == nil {
		t.Fatal("expected provider ID to remain unchanged")
	}

	if *providerID != fixture.ProviderID {
		t.Fatalf(
			"expected provider ID %d, got %d",
			fixture.ProviderID,
			*providerID,
		)
	}

	if externalOrderID == "" {
		t.Fatal("expected external order ID to remain unchanged")
	}
}

func TestExecuteOrderRefund_RollsBackAllChangesWhenFinalStateUpdateFails(t *testing.T) {

	checkoutRepo, pool := newTestCheckoutRepository(t)

	refundAmount := decimal.NewFromInt(30)

	fixture := createRefundFixture(
		t,
		pool,
		orderentity.OrderStatusProcessing,
		refundAmount,
	)

	beforeBalance := getWalletBalance(
		t,
		pool,
		fixture.WalletID,
	)

	_, providerBefore, externalBefore := getOrderState(
		t,
		pool,
		fixture.OrderID,
	)

	if providerBefore == nil {
		t.Fatal("expected provider ID before refund")
	}

	installFailOrderStatusTransition(
		t,
		pool,
		fixture.OrderID,
	)

	_, err := checkoutRepo.ExecuteOrderRefund(
		context.Background(),
		checkoutparams.RefundOrderRequest{
			OrderID: fixture.OrderID,
			Reason:  "provider_failed",
		},
	)

	if err == nil {
		t.Fatal("expected refund transaction to fail")
	}

	afterBalance := getWalletBalance(
		t,
		pool,
		fixture.WalletID,
	)

	if !afterBalance.Equal(beforeBalance) {
		t.Fatalf(
			"expected wallet balance to rollback to %s, got %s",
			beforeBalance,
			afterBalance,
		)
	}

	txCount := getRefundTransactionCount(
		t,
		pool,
		fixture.OrderID,
	)

	if txCount != 0 {
		t.Fatalf(
			"expected refund transaction to be rolled back, got %d",
			txCount,
		)
	}

	status, providerAfter, externalAfter := getOrderState(
		t,
		pool,
		fixture.OrderID,
	)

	if status != orderentity.OrderStatusProcessing {
		t.Fatalf(
			"expected order status PROCESSING after rollback, got %s",
			status,
		)
	}

	if providerAfter == nil {
		t.Fatal("expected provider ID to remain after rollback")
	}

	if *providerAfter != *providerBefore {
		t.Fatalf(
			"expected provider ID %d after rollback, got %d",
			*providerBefore,
			*providerAfter,
		)
	}

	if externalAfter != externalBefore {
		t.Fatalf(
			"expected external order ID %s after rollback, got %s",
			externalBefore,
			externalAfter,
		)
	}
}
