package paymentflowtesting

import (
	"context"
	"testing"

	"telegram-service-platform/entity"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/repository/postgres"
	"telegram-service-platform/repository/postgrespayment"
	"telegram-service-platform/service/paymentservice"

	"github.com/shopspring/decimal"
)

func TestStartPaymentFlow_NewPayment(t *testing.T) {

	pool := newTestPool(t)

	var (
		paymentID uint64
		orderID   uint64
	)

	t.Cleanup(func() {
		cleanupPaymentTestData(t, pool, paymentID, orderID)
	})

	order := createTestOrder(
		t,
		pool,
		orderentity.OrderStatusPending,
	)
	orderID = order.ID
	provider := &flowTestProvider{
		createResponse: paymentproviderparams.CreateResponse{
			ExternalID: "AUTH-START-123",
			PaymentURL: "https://provider.test/pay/AUTH-START-123",
		},
	}

	transactionProvider := postgres.NewTransactionProvider(pool)

	paymentRepo := postgrespayment.NewWithExecutor(
		pool,
		transactionProvider,
	)

	service := paymentservice.New(
		paymentRepo,
		paymentRepo,
		provider,
		nil,
	)

	req := paymentparams.StartPaymentRequest{
		OrderID:        order.ID,
		UserID:         1,
		Method:         paymententity.PaymentMethodZarinpal,
		Amount:         entity.Amount(decimal.NewFromInt(100000)),
		Currency:       entity.CurrencyTOMAN,
		IdempotencyKey: "start-payment-flow-001",
		CallbackURL:    "https://example.com/payments/zarinpal/callback",
		Description:    "Start payment integration test",
	}

	var beforeCount int

	err := pool.QueryRow(
		context.Background(),
		`
			SELECT COUNT(*)
			FROM payments
			WHERE idempotency_key = $1
		`,
		req.IdempotencyKey,
	).Scan(&beforeCount)
	if err != nil {
		t.Fatalf(
			"check initial payment state: %v",
			err,
		)
	}

	if beforeCount != 0 {
		t.Fatalf(
			"expected no payment before StartPayment, got %d",
			beforeCount,
		)
	}

	response, sErr := service.StartPayment(
		context.Background(),
		req,
	)
	if sErr != nil {
		t.Fatalf(
			"unexpected StartPayment error: %v",
			sErr,
		)
	}

	paymentID = response.PaymentID

	if response.PaymentID == 0 {
		t.Fatal("expected payment ID, got 0")
	}

	if response.Status != paymententity.PaymentStatusPending {
		t.Fatalf("expected response status PENDING, got %s", response.Status)
	}

	if response.ExternalID != "AUTH-START-123" {
		t.Fatalf("expected external ID AUTH-START-123, got %s", response.ExternalID)
	}

	if response.PaymentURL != "https://provider.test/pay/AUTH-START-123" {
		t.Fatalf("unexpected payment URL: %s", response.PaymentURL)
	}

	if provider.createCalls != 1 {
		t.Fatalf("expected provider Create to be called once, got %d", provider.createCalls)
	}

	payment := readPayment(t, pool, paymentID)

	if payment.OrderID != order.ID {
		t.Fatalf(
			"expected payment order id %d, got %d",
			order.ID,
			payment.OrderID,
		)
	}

	if payment.UserID != req.UserID {
		t.Fatalf(
			"expected payment user id %d, got %d",
			req.UserID,
			payment.UserID,
		)
	}

	if payment.Method != req.Method {
		t.Fatalf(
			"expected payment method %s, got %s",
			req.Method,
			payment.Method,
		)
	}

	if !payment.Amount.Equal(req.Amount) {
		t.Fatalf(
			"expected payment amount %s, got %s",
			req.Amount,
			payment.Amount,
		)
	}

	if payment.Currency != req.Currency {
		t.Fatalf(
			"expected payment currency %s, got %s",
			req.Currency,
			payment.Currency,
		)
	}

	if payment.Status != paymententity.PaymentStatusPending {
		t.Fatalf(
			"expected persisted payment status PENDING, got %s",
			payment.Status,
		)
	}

	if payment.ExternalID != "AUTH-START-123" {
		t.Fatalf(
			"expected persisted external ID AUTH-START-123, got %s",
			payment.ExternalID,
		)
	}

	if payment.PaymentURL !=
		"https://provider.test/pay/AUTH-START-123" {
		t.Fatalf(
			"unexpected persisted payment URL: %s",
			payment.PaymentURL,
		)
	}

	if payment.IdempotencyKey != req.IdempotencyKey {
		t.Fatalf(
			"expected idempotency key %s, got %s",
			req.IdempotencyKey,
			payment.IdempotencyKey,
		)
	}
}
