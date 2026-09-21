package checkouttesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/params/checkoutparams"
	"telegram-service-platform/service/checkoutservice"
	"telegram-service-platform/service/notificationservice"
	"testing"

	"github.com/shopspring/decimal"
)

func newCheckoutService(
	transactionRepo *fakeTransactionRepository,
	notificationSvc *notificationservice.Service,
) *checkoutservice.Service {
	return checkoutservice.New(
		nil,
		nil,
		nil,
		nil,
		notificationSvc,
		transactionRepo,
		nil,
		nil,
		checkoutservice.Config{},
	)
}

func TestService_RefundOrder_Success(t *testing.T) {

	transactionRepo := &fakeTransactionRepository{
		response: &checkoutparams.RefundOrderResponse{
			OrderID:      42,
			UserID:       100,
			WalletTxID:   500,
			RefundAmount: entity.Amount(decimal.NewFromInt(100)),
		},
	}

	notificationRepo := &fakeNotificationRepository{}
	notificationRedis := &fakeNotificationRedis{}

	notificationSvc := newNotificationService(
		notificationRepo,
		notificationRedis,
	)

	service := newCheckoutService(
		transactionRepo,
		notificationSvc,
	)

	req := checkoutparams.RefundOrderRequest{
		OrderID: 42,
		Reason:  "provider_failed",
	}

	err := service.RefundOrder(
		context.Background(),
		req,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(transactionRepo.calls) != 1 {
		t.Fatalf(
			"expected one refund transaction call, got %d",
			len(transactionRepo.calls),
		)
	}

	refundReq := transactionRepo.calls[0]

	if refundReq.OrderID != 42 {
		t.Fatalf(
			"expected order ID 42, got %d",
			refundReq.OrderID,
		)
	}

	if refundReq.Reason != "provider_failed" {
		t.Fatalf(
			"expected reason provider_failed, got %s",
			refundReq.Reason,
		)
	}

	if len(notificationRepo.created) != 1 {
		t.Fatalf(
			"expected one notification, got %d",
			len(notificationRepo.created),
		)
	}

	notification := notificationRepo.created[0]

	if notification.UserID != 100 {
		t.Fatalf(
			"expected notification user ID 100, got %d",
			notification.UserID,
		)
	}

	if notification.Type != notificationentity.NotificationTypeOrderFailed {
		t.Fatalf(
			"expected order failed notification, got %s",
			notification.Type,
		)
	}

	if len(notificationRedis.enqueued) != 1 {
		t.Fatalf(
			"expected one notification enqueue, got %d",
			len(notificationRedis.enqueued),
		)
	}
}

func TestService_RefundOrder_TransactionFailureDoesNotNotify(t *testing.T) {

	transactionRepo := &fakeTransactionRepository{
		err: errors.New("refund transaction failed"),
	}

	notificationRepo := &fakeNotificationRepository{}
	notificationRedis := &fakeNotificationRedis{}

	notificationSvc := newNotificationService(
		notificationRepo,
		notificationRedis,
	)

	service := newCheckoutService(
		transactionRepo,
		notificationSvc,
	)

	err := service.RefundOrder(
		context.Background(),
		checkoutparams.RefundOrderRequest{
			OrderID: 42,
			Reason:  "provider_failed",
		},
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if len(transactionRepo.calls) != 1 {
		t.Fatalf(
			"expected one refund call, got %d",
			len(transactionRepo.calls),
		)
	}

	if len(notificationRepo.created) != 0 {
		t.Fatalf(
			"expected no notification, got %d",
			len(notificationRepo.created),
		)
	}

	if len(notificationRedis.enqueued) != 0 {
		t.Fatalf(
			"expected no notification enqueue, got %d",
			len(notificationRedis.enqueued),
		)
	}
}
