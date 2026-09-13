package paymentflowtesting

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/pkg/richerror"
	"testing"
)

func TestZarinpalCallbackFlow_Success(t *testing.T) {
	provider := &flowTestProvider{
		verifyResponse: paymentproviderparams.VerifyResponse{
			Status:      paymententity.PaymentStatusSuccess,
			ReferenceID: "REF-123",
		},
	}

	validator := &flowTestValidator{}

	echoServer, pool := newZarinpalCallbackFlow(t, provider, validator)

	server := httptest.NewServer(echoServer)
	defer server.Close()

	payment, order := createTestCase(
		t,
		pool,
		paymententity.PaymentStatusPending,
		orderentity.OrderStatusPending,
	)

	setPaymentExternalID(t, pool, payment.ID, "AUTH-123")

	response, err := http.Get(
		server.URL +
			"/payments/zarinpal/callback" +
			"?status=OK&authority=AUTH-123",
	)
	if err != nil {
		t.Fatalf("send callback: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}

	if validator.calls != 1 {
		t.Fatalf("expected validator to be called once, got %d", validator.calls)
	}

	if provider.verifyCalls != 1 {
		t.Fatalf("expected provider Verify to be called once, got %d", provider.verifyCalls)
	}

	if provider.lastVerifyRequest.ExternalID != "AUTH-123" {
		t.Fatalf("expected external ID AUTH-123, got %s", provider.lastVerifyRequest.ExternalID)
	}

	paymentStatus, referenceID := readPaymentState(t, pool, payment.ID)

	if paymentStatus != paymententity.PaymentStatusSuccess {
		t.Fatalf("expected payment SUCCESS, got %s", paymentStatus)
	}

	if referenceID != "REF-123" {
		t.Fatalf("expected reference REF-123, got %s", referenceID)
	}

	orderStatus := readOrderStatus(t, pool, order.ID)

	if orderStatus != orderentity.OrderStatusPaid {
		t.Fatalf("expected order PAID, got %s", orderStatus)
	}
}

func TestZarinpalCallbackFlow_Timeout(
	t *testing.T,
) {
	provider := &flowTestProvider{
		verifyErr: richerror.New(
			"fake.zarinpal.verify",
			context.DeadlineExceeded,
		).WithCode(
			richerror.CodePaymentProviderTimeout,
		),
	}

	validator := &flowTestValidator{}

	echoServer, pool := newZarinpalCallbackFlow(t, provider, validator)

	server := httptest.NewServer(echoServer)
	defer server.Close()

	payment, order := createTestCase(
		t,
		pool,
		paymententity.PaymentStatusPending,
		orderentity.OrderStatusPending,
	)

	setPaymentExternalID(t, pool, payment.ID, "AUTH-TIMEOUT")

	response, err := http.Get(
		server.URL +
			"/payments/zarinpal/callback" +
			"?status=OK&authority=AUTH-TIMEOUT",
	)
	if err != nil {
		t.Fatalf("send callback: %v", err)
	}

	defer response.Body.Close()

	if provider.verifyCalls != 1 {
		t.Fatalf("expected provider Verify once, got %d", provider.verifyCalls)
	}

	status, _ := readPaymentState(t, pool, payment.ID)

	if status != paymententity.PaymentStatusUnknown {
		t.Fatalf("expected payment UNKNOWN, got %s", status)
	}

	orderStatus := readOrderStatus(t, pool, order.ID)

	if orderStatus != orderentity.OrderStatusPending {
		t.Fatalf("expected order PENDING, got %s", orderStatus)
	}
}

func TestZarinpalCallbackFlow_Rejected(
	t *testing.T,
) {
	provider := &flowTestProvider{
		verifyErr: richerror.New(
			"fake.zarinpal.verify",
			errors.New("payment rejected"),
		).WithCode(
			richerror.CodePaymentProviderRejected,
		),
	}

	validator := &flowTestValidator{}

	echoServer, pool := newZarinpalCallbackFlow(t, provider, validator)

	server := httptest.NewServer(echoServer)
	defer server.Close()

	payment, order := createTestCase(
		t,
		pool,
		paymententity.PaymentStatusFailed,
		orderentity.OrderStatusPending,
	)

	setPaymentExternalID(
		t,
		pool,
		payment.ID,
		"AUTH-REJECTED",
	)

	response, err := http.Get(
		server.URL +
			"/payments/zarinpal/callback" +
			"?status=OK&authority=AUTH-REJECTED",
	)
	if err != nil {
		t.Fatalf("send callback: %v", err)
	}

	defer response.Body.Close()

	status, _ := readPaymentState(t, pool, payment.ID)

	if status != paymententity.PaymentStatusFailed {
		t.Fatalf("expected payment FAILED, got %s", status)
	}

	orderStatus := readOrderStatus(t, pool, order.ID)

	if orderStatus != orderentity.OrderStatusPending {
		t.Fatalf("expected order PENDING, got %s", orderStatus)
	}
}

func TestZarinpalCallbackFlow_Rollback(
	t *testing.T,
) {
	provider := &flowTestProvider{
		verifyResponse: paymentproviderparams.VerifyResponse{
			Status:      paymententity.PaymentStatusSuccess,
			ReferenceID: "REF-ROLLBACK",
		},
	}

	validator := &flowTestValidator{}

	echoServer, pool := newZarinpalCallbackFlow(
		t,
		provider,
		validator,
	)

	server := httptest.NewServer(echoServer)
	defer server.Close()

	payment, order := createTestCase(
		t,
		pool,
		paymententity.PaymentStatusPending,
		orderentity.OrderStatusCanceled,
	)

	setPaymentExternalID(
		t,
		pool,
		payment.ID,
		"AUTH-ROLLBACK",
	)

	_, err := http.Get(
		server.URL +
			"/payments/zarinpal/callback" +
			"?status=OK&authority=AUTH-ROLLBACK",
	)
	if err != nil {
		t.Fatalf("send callback: %v", err)
	}

	paymentStatus, referenceID := readPaymentState(t, pool, payment.ID)

	if paymentStatus != paymententity.PaymentStatusPending {
		t.Fatalf("expected payment PENDING after rollback, got %s", paymentStatus)
	}

	if referenceID != "" {
		t.Fatalf("expected empty provider reference after rollback, got %s", referenceID)
	}

	orderStatus := readOrderStatus(t, pool, order.ID)

	if orderStatus != orderentity.OrderStatusCanceled {
		t.Fatalf("expected order CANCELED, got %s", orderStatus)
	}
}

func TestZarinpalCallbackFlow_DuplicateCallback(
	t *testing.T,
) {
	provider := &flowTestProvider{
		verifyResponse: paymentproviderparams.VerifyResponse{
			Status:      paymententity.PaymentStatusSuccess,
			ReferenceID: "REF-DUP",
		},
	}

	validator := &flowTestValidator{}

	echoServer, pool := newZarinpalCallbackFlow(t, provider, validator)

	server := httptest.NewServer(echoServer)
	defer server.Close()

	payment, order := createTestCase(
		t,
		pool,
		paymententity.PaymentStatusPending,
		orderentity.OrderStatusPending,
	)

	setPaymentExternalID(
		t,
		pool,
		payment.ID,
		"AUTH-DUP",
	)

	callbackURL :=
		server.URL +
			"/payments/zarinpal/callback" +
			"?status=OK&authority=AUTH-DUP"

	firstResponse, err := http.Get(callbackURL)
	if err != nil {
		t.Fatalf("first callback: %v", err)
	}

	firstResponse.Body.Close()

	if firstResponse.StatusCode != http.StatusOK {
		t.Fatalf("expected first callback 200, got %d", firstResponse.StatusCode)
	}

	secondResponse, gErr := http.Get(callbackURL)
	if gErr != nil {
		t.Fatalf("second callback: %v", gErr)
	}

	secondResponse.Body.Close()

	if secondResponse.StatusCode != http.StatusOK {
		t.Fatalf("expected second callback 200, got %d", secondResponse.StatusCode)
	}

	if provider.verifyCalls != 1 {
		t.Fatalf("expected provider Verify to be called once, got %d", provider.verifyCalls)
	}

	status, referenceID := readPaymentState(t, pool, payment.ID)

	if status != paymententity.PaymentStatusSuccess {
		t.Fatalf("expected SUCCESS, got %s", status)
	}

	if referenceID != "REF-DUP" {
		t.Fatalf("expected REF-DUP, got %s", referenceID)
	}

	orderStatus := readOrderStatus(t, pool, order.ID)

	if orderStatus != orderentity.OrderStatusPaid {
		t.Fatalf("expected PAID, got %s", orderStatus)
	}
}
