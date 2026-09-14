package paymenthandlertesting

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"telegram-service-platform/config"
	"telegram-service-platform/delivery/httpserver/paymenthandler"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/service/authservice"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestHandler_StartPayment_Success(t *testing.T) {
	e := echo.New()

	flow := &fakePaymentFlow{
		startPaymentResponse: &paymentparams.InitiateResponse{
			PaymentID:  100,
			Status:     paymententity.PaymentStatusPending,
			ExternalID: "AUTH-100",
			PaymentURL: "https://provider.test/pay/AUTH-100",
		},
	}

	validator := &fakePaymentValidator{}

	handler := paymenthandler.New(flow, validator, nil)

	body := map[string]any{
		"order_id":        10,
		"method":          string(paymententity.PaymentMethodZarinpal),
		"amount":          "100000",
		"currency":        string(entity.CurrencyTOMAN),
		"idempotency_key": "idem-100",
		"callback_url":    "https://example.com/payments/zarinpal/callback",
		"description":     "Order #10",
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/payments",
		bytes.NewReader(bodyBytes),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	ctx.Set(
		config.AuthMiddlewareContextKey,
		&authservice.Claims{
			UserId: 20,
		},
	)

	err = handler.StartPaymentHandler(ctx)
	if err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !validator.validateStartPaymentCalled {
		t.Fatal("expected ValidateStartPayment to be called")
	}

	if !flow.startPaymentCalled {
		t.Fatal("expected StartPayment to be called")
	}

	if flow.startPaymentReq.UserID != 20 {
		t.Fatalf("expected user id 20 from claims, got %d", flow.startPaymentReq.UserID)
	}

	if flow.startPaymentReq.OrderID != 10 {
		t.Fatalf("expected order id 10, got %d", flow.startPaymentReq.OrderID)
	}

	if flow.startPaymentReq.Method !=
		paymententity.PaymentMethodZarinpal {
		t.Fatalf(
			"expected method %s, got %s",
			paymententity.PaymentMethodZarinpal,
			flow.startPaymentReq.Method,
		)
	}

	if flow.startPaymentReq.Currency != entity.CurrencyTOMAN {
		t.Fatalf(
			"expected currency %s, got %s",
			entity.CurrencyTOMAN,
			flow.startPaymentReq.Currency,
		)
	}

	if flow.startPaymentReq.IdempotencyKey != "idem-100" {
		t.Fatalf(
			"expected idempotency key idem-100, got %s",
			flow.startPaymentReq.IdempotencyKey,
		)
	}

	if flow.startPaymentReq.CallbackURL !=
		"https://example.com/payments/zarinpal/callback" {
		t.Fatalf(
			"unexpected callback URL: %s",
			flow.startPaymentReq.CallbackURL,
		)
	}

	if flow.startPaymentReq.Description != "Order #10" {
		t.Fatalf(
			"unexpected description: %s",
			flow.startPaymentReq.Description,
		)
	}

	var response paymentparams.InitiateResponse

	if err := json.Unmarshal(
		rec.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatalf(
			"decode response: %v",
			err,
		)
	}

	if response.PaymentID != 100 {
		t.Fatalf(
			"expected payment id 100, got %d",
			response.PaymentID,
		)
	}

	if response.Status != paymententity.PaymentStatusPending {
		t.Fatalf(
			"expected status PENDING, got %s",
			response.Status,
		)
	}

	if response.ExternalID != "AUTH-100" {
		t.Fatalf(
			"expected external id AUTH-100, got %s",
			response.ExternalID,
		)
	}

	if response.PaymentURL !=
		"https://provider.test/pay/AUTH-100" {
		t.Fatalf(
			"unexpected payment URL: %s",
			response.PaymentURL,
		)
	}
}
