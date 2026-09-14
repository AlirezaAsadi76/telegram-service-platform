package paymentflowtesting

import (
	"context"
	"errors"
	"telegram-service-platform/pkg/richerror"
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

func TestStartPaymentFlow_ExistingPaymentByIdempotencyKey(t *testing.T) {

	pool := newTestPool(t)

	var (
		paymentID uint64
		orderID   uint64
	)

	t.Cleanup(func() {
		cleanupPaymentTestData(
			t,
			pool,
			paymentID,
			orderID,
		)
	})

	order := createTestOrder(
		t,
		pool,
		orderentity.OrderStatusPending,
	)

	orderID = order.ID

	transactionProvider := postgres.NewTransactionProvider(pool)

	paymentRepo := postgrespayment.NewWithExecutor(
		pool,
		transactionProvider,
	)

	provider := &flowTestProvider{
		createResponse: paymentproviderparams.CreateResponse{
			ExternalID: "SHOULD-NOT-BE-USED",
			PaymentURL: "https://provider.test/should-not-be-used",
		},
	}

	service := paymentservice.New(
		paymentRepo,
		paymentRepo,
		provider,
		nil,
	)

	const idempotencyKey = "start-payment-idempotency-001"

	existingPayment := &paymententity.Payment{
		OrderID:        order.ID,
		UserID:         1,
		Method:         paymententity.PaymentMethodZarinpal,
		Amount:         entity.Amount(decimal.NewFromInt(100000)),
		Currency:       entity.CurrencyTOMAN,
		Status:         paymententity.PaymentStatusPending,
		ExternalID:     "AUTH-EXISTING-001",
		PaymentURL:     "https://provider.test/pay/AUTH-EXISTING-001",
		IdempotencyKey: idempotencyKey,
	}

	err := paymentRepo.Create(
		context.Background(),
		existingPayment,
	)
	if err != nil {
		t.Fatalf("create existing payment: %v", err)
	}

	paymentID = existingPayment.ID

	req := paymentparams.StartPaymentRequest{
		OrderID:        order.ID,
		UserID:         existingPayment.UserID,
		Method:         existingPayment.Method,
		Amount:         existingPayment.Amount,
		Currency:       existingPayment.Currency,
		IdempotencyKey: idempotencyKey,
		CallbackURL:    "https://example.com/payments/zarinpal/callback",
		Description:    "Idempotency test",
	}

	var initialCount int

	err = pool.QueryRow(
		context.Background(),
		`
			SELECT COUNT(*)
			FROM payments
			WHERE idempotency_key = $1
		`,
		idempotencyKey,
	).Scan(&initialCount)
	if err != nil {
		t.Fatalf("check initial payment count: %v", err)
	}

	if initialCount != 1 {
		t.Fatalf("expected exactly one existing payment, got %d", initialCount)
	}

	response, sErr := service.StartPayment(
		context.Background(),
		req,
	)
	if sErr != nil {
		t.Fatalf("unexpected StartPayment error: %v", sErr)
	}

	if response == nil {
		t.Fatal("expected response, got nil")
	}

	if response.PaymentID != existingPayment.ID {
		t.Fatalf(
			"expected existing payment id %d, got %d",
			existingPayment.ID,
			response.PaymentID,
		)
	}

	if response.Status != paymententity.PaymentStatusPending {
		t.Fatalf(
			"expected response status PENDING, got %s",
			response.Status,
		)
	}

	if response.ExternalID != existingPayment.ExternalID {
		t.Fatalf(
			"expected external id %s, got %s",
			existingPayment.ExternalID,
			response.ExternalID,
		)
	}

	if response.PaymentURL != existingPayment.PaymentURL {
		t.Fatalf(
			"expected payment url %s, got %s",
			existingPayment.PaymentURL,
			response.PaymentURL,
		)
	}

	if provider.createCalls != 0 {
		t.Fatalf(
			"expected provider Create not to be called, got %d calls",
			provider.createCalls,
		)
	}

	finalCount := countPaymentByIdempotencyKey(t, pool, idempotencyKey)

	if finalCount != 1 {
		t.Fatalf(
			"expected exactly one payment after repeated request, got %d",
			finalCount,
		)
	}

	payment := readPayment(
		t,
		pool,
		existingPayment.ID,
	)

	if payment.ID != existingPayment.ID {
		t.Fatalf(
			"expected payment id %d, got %d",
			existingPayment.ID,
			payment.ID,
		)
	}

	if payment.Status != paymententity.PaymentStatusPending {
		t.Fatalf(
			"expected persisted status PENDING, got %s",
			payment.Status,
		)
	}

	if payment.ExternalID != existingPayment.ExternalID {
		t.Fatalf(
			"expected persisted external id %s, got %s",
			existingPayment.ExternalID,
			payment.ExternalID,
		)
	}

	if payment.PaymentURL != existingPayment.PaymentURL {
		t.Fatalf(
			"expected persisted payment url %s, got %s",
			existingPayment.PaymentURL,
			payment.PaymentURL,
		)
	}

	if payment.IdempotencyKey != idempotencyKey {
		t.Fatalf(
			"expected persisted idempotency key %s, got %s",
			idempotencyKey,
			payment.IdempotencyKey,
		)
	}
}

func TestStartPaymentFlow_ProviderFailureMarksPaymentUnknown(t *testing.T) {

	pool := newTestPool(t)

	var (
		paymentID uint64
		orderID   uint64
	)

	t.Cleanup(func() {
		cleanupPaymentTestData(
			t,
			pool,
			paymentID,
			orderID,
		)
	})

	order := createTestOrder(
		t,
		pool,
		orderentity.OrderStatusPending,
	)

	orderID = order.ID

	provider := &flowTestProvider{
		createErr: errors.New("provider connection failed"),
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
		IdempotencyKey: "start-payment-provider-failure-001",
		CallbackURL:    "https://example.com/payments/zarinpal/callback",
		Description:    "Provider failure integration test",
	}

	beforeCount := countPaymentByIdempotencyKey(t, pool, req.IdempotencyKey)

	if beforeCount != 0 {
		t.Fatalf(
			"expected no payment before StartPayment, got %d",
			beforeCount,
		)
	}

	response, err := service.StartPayment(
		context.Background(),
		req,
	)

	if err == nil {
		t.Fatal("expected StartPayment to return an error")
	}

	if response != nil {
		t.Fatal(
			"expected response to be nil when provider creation fails",
		)
	}

	if provider.createCalls != 1 {
		t.Fatalf(
			"expected provider Create to be called once, got %d",
			provider.createCalls,
		)
	}

	paymentCount := countPaymentByIdempotencyKey(t, pool, req.IdempotencyKey)

	if paymentCount != 1 {
		t.Fatalf(
			"expected one payment after provider failure, got %d",
			paymentCount,
		)
	}

	payment := readPayment(
		t,
		pool,
		getPaymentIDByIdempotencyKey(t, pool, req.IdempotencyKey),
	)

	paymentID = payment.ID

	if payment.OrderID != req.OrderID {
		t.Fatalf(
			"expected payment order id %d, got %d",
			req.OrderID,
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

	if payment.Status != paymententity.PaymentStatusUnknown {
		t.Fatalf(
			"expected persisted payment UNKNOWN, got %s",
			payment.Status,
		)
	}

	if payment.ExternalID != "" {
		t.Fatalf(
			"expected empty external id, got %s",
			payment.ExternalID,
		)
	}

	if payment.PaymentURL != "" {
		t.Fatalf(
			"expected empty payment url, got %s",
			payment.PaymentURL,
		)
	}

	orderStatus := readOrderStatus(t, pool, orderID)

	if orderStatus != orderentity.OrderStatusPending {
		t.Fatalf(
			"expected order status PENDING, got %s",
			orderStatus,
		)
	}
}

func TestStartPaymentFlow_ProviderRejectedMarksPaymentFailed(t *testing.T) {

	pool := newTestPool(t)

	var (
		paymentID uint64
		orderID   uint64
	)

	t.Cleanup(func() {
		cleanupPaymentTestData(
			t,
			pool,
			paymentID,
			orderID,
		)
	})

	order := createTestOrder(
		t,
		pool,
		orderentity.OrderStatusPending,
	)

	orderID = order.ID

	rejectionErr := richerror.New(
		"fakeprovider.create",
		errors.New("provider rejected payment"),
	).
		WithKind(richerror.KindExternalAPI).
		WithCode(richerror.CodePaymentProviderRejected)

	provider := &flowTestProvider{
		createErr: rejectionErr,
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
		IdempotencyKey: "start-payment-provider-rejected-001",
		CallbackURL:    "https://example.com/payments/zarinpal/callback",
		Description:    "Provider rejection integration test",
	}

	initialCount := countPaymentByIdempotencyKey(t, pool, req.IdempotencyKey)
	if initialCount != 0 {
		t.Fatalf(
			"expected no payment before StartPayment, got %d",
			initialCount,
		)
	}

	response, err := service.StartPayment(
		context.Background(),
		req,
	)

	if err == nil {
		t.Fatal("expected StartPayment to return an error")
	}

	if response != nil {
		t.Fatal(
			"expected response to be nil when provider rejects payment",
		)
	}

	if !richerror.IsCode(err, richerror.CodePaymentProviderRejected) {
		t.Fatalf(
			"expected error code %s, got %s",
			richerror.CodePaymentProviderRejected,
			richerror.New("", err).Code(),
		)
	}

	if provider.createCalls != 1 {
		t.Fatalf(
			"expected provider Create to be called once, got %d",
			provider.createCalls,
		)
	}

	paymentID = getPaymentIDByIdempotencyKey(
		t,
		pool,
		req.IdempotencyKey,
	)

	payment := readPayment(
		t,
		pool,
		paymentID,
	)

	if payment.Status != paymententity.PaymentStatusFailed {
		t.Fatalf(
			"expected persisted payment FAILED, got %s",
			payment.Status,
		)
	}

	if payment.ExternalID != "" {
		t.Fatalf(
			"expected empty external id, got %s",
			payment.ExternalID,
		)
	}

	if payment.PaymentURL != "" {
		t.Fatalf(
			"expected empty payment url, got %s",
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

	orderStatus := readOrderStatus(t, pool, orderID)

	if orderStatus != orderentity.OrderStatusPending {
		t.Fatalf(
			"expected order status PENDING, got %s",
			orderStatus,
		)
	}
}
