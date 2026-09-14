package paymenthandlertesting

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"telegram-service-platform/config"
	"telegram-service-platform/delivery/httpserver/paymenthandler"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
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

func TestHandler_StartPayment_ValidationError(t *testing.T) {

	e := echo.New()

	flow := &fakePaymentFlow{}

	validationErr := richerror.New(
		"validator.validateStartPayment",
		nil,
	).
		WithMessage(msgerror.InvalidInput).
		WithKind(richerror.KindInvalid)

	validator := &fakePaymentValidator{
		fieldErrs: map[string]string{
			"Amount": "is required",
		},
		err: validationErr,
	}

	handler := paymenthandler.New(
		flow,
		validator,
		nil,
	)

	body := map[string]any{
		"order_id":        10,
		"method":          string(paymententity.PaymentMethodZarinpal),
		"amount":          "100000",
		"currency":        string(entity.CurrencyTOMAN),
		"idempotency_key": "idem-validation-001",
		"callback_url":    "https://example.com/payments/zarinpal/callback",
		"description":     "Validation error test",
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf(
			"marshal request body: %v",
			err,
		)
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

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, rec.Code)
	}

	if !validator.validateStartPaymentCalled {
		t.Fatal("expected ValidateStartPayment to be called")
	}

	if flow.startPaymentCalled {
		t.Fatal("expected StartPayment not to be called")
	}

	var response params.ValidationErrorResponse

	if err := json.Unmarshal(
		rec.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatalf(
			"decode validation error response: %v",
			err,
		)
	}

	if response.Message != msgerror.InvalidInput {
		t.Fatalf("expected message %q, got %q", msgerror.InvalidInput, response.Message)
	}

	if response.FieldErrors["Amount"] != "is required" {
		t.Fatalf("expected Amount field error %q, got %q", "is required", response.FieldErrors["Amount"])
	}
}

func TestHandler_StartPayment_ServiceConflict(t *testing.T) {

	e := echo.New()

	flow := &fakePaymentFlow{
		startPaymentError: richerror.New(
			"paymentservice.start_payment",
			errors.New("internal database or business details"),
		).
			WithKind(richerror.KindConflict).
			WithCode(richerror.CodePaymentInvalidState).
			WithMessage("payment conflict"),
	}

	validator := &fakePaymentValidator{}

	handler := paymenthandler.New(
		flow,
		validator,
		nil,
	)

	body := map[string]any{
		"order_id":        10,
		"method":          string(paymententity.PaymentMethodZarinpal),
		"amount":          "100000",
		"currency":        string(entity.CurrencyTOMAN),
		"idempotency_key": "service-conflict-001",
		"callback_url":    "https://example.com/payments/zarinpal/callback",
		"description":     "Service conflict test",
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

	if err == nil {
		t.Fatal("expected handler to return an HTTP error")
	}

	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)

	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T", err)
	}

	if httpErr.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, httpErr.Code)
	}

	if httpErr.Message != "payment conflict" {
		t.Fatalf("expected message %q, got %v", "payment conflict", httpErr.Message)
	}

	if !flow.startPaymentCalled {
		t.Fatal("expected StartPayment to be called")
	}

	if validator.validateStartPaymentCalled == false {
		t.Fatal("expected ValidateStartPayment to be called")
	}
}

func TestHandler_StartPayment_InternalErrorDoesNotLeakDetails(t *testing.T) {

	e := echo.New()

	const internalMessage = "pq: connection refused at postgres.internal:5432"

	flow := &fakePaymentFlow{
		startPaymentError: richerror.New(
			"paymentservice.start_payment",
			nil,
		).
			WithKind(richerror.KindUnexpected).
			WithMessage(internalMessage),
	}

	validator := &fakePaymentValidator{}

	handler := paymenthandler.New(
		flow,
		validator,
		nil,
	)

	body := map[string]any{
		"order_id":        10,
		"method":          string(paymententity.PaymentMethodZarinpal),
		"amount":          "100000",
		"currency":        string(entity.CurrencyTOMAN),
		"idempotency_key": "internal-error-001",
		"callback_url":    "https://example.com/payments/zarinpal/callback",
		"description":     "Internal error test",
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

	if err == nil {
		t.Fatal("expected handler to return an HTTP error")
	}

	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T", err)
	}

	if httpErr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, httpErr.Code)
	}

	if httpErr.Message == internalMessage {
		t.Fatal(
			"technical internal error leaked to client",
		)
	}

	if !validator.validateStartPaymentCalled {
		t.Fatal("expected ValidateStartPayment to be called")
	}

	if !flow.startPaymentCalled {
		t.Fatal(
			"expected StartPayment to be called",
		)
	}
}

func TestHandler_StartPayment_MissingClaims(t *testing.T) {

	e := echo.New()

	flow := &fakePaymentFlow{}

	validator := &fakePaymentValidator{}

	handler := paymenthandler.New(
		flow,
		validator,
		nil,
	)

	body := map[string]any{
		"order_id":        10,
		"method":          string(paymententity.PaymentMethodZarinpal),
		"amount":          "100000",
		"currency":        string(entity.CurrencyTOMAN),
		"idempotency_key": "auth-error-001",
		"callback_url":    "https://example.com/payments/zarinpal/callback",
		"description":     "Authentication test",
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

	err = handler.StartPaymentHandler(ctx)

	if err == nil {
		t.Fatal("expected handler to return unauthorized error")
	}

	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T", err)
	}

	if httpErr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, httpErr.Code)
	}

	if httpErr.Message != "unauthorized" {
		t.Fatalf("expected message %q, got %v", "unauthorized", httpErr.Message)
	}

	if !validator.validateStartPaymentCalled {
		t.Fatal("expected validation to run before authentication check")
	}

	if flow.startPaymentCalled {
		t.Fatal("expected StartPayment not to be called")
	}

	if ctx.Get(config.AuthMiddlewareContextKey) != nil {
		t.Fatal("expected authentication claims to be absent")
	}
}

func TestHandler_StartPayment_InvalidClaimsType(t *testing.T) {
	e := echo.New()

	flow := &fakePaymentFlow{}
	validator := &fakePaymentValidator{}

	handler := paymenthandler.New(
		flow,
		validator,
		nil,
	)

	body := map[string]any{
		"order_id":        10,
		"method":          string(paymententity.PaymentMethodZarinpal),
		"amount":          "100000",
		"currency":        string(entity.CurrencyTOMAN),
		"idempotency_key": "invalid-claims-type-001",
		"callback_url":    "https://example.com/payments/zarinpal/callback",
		"description":     "Invalid claims type test",
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

	ctx.Set(config.AuthMiddlewareContextKey, "invalid-claims")

	err = handler.StartPaymentHandler(ctx)

	if err == nil {
		t.Fatal("expected handler to return unauthorized error")
	}

	var httpErr *echo.HTTPError
	ok := errors.As(err, &httpErr)
	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T", err)
	}

	if httpErr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, httpErr.Code)
	}

	if httpErr.Message != "unauthorized" {
		t.Fatalf("expected message %q, got %v", "unauthorized", httpErr.Message)
	}

	if !validator.validateStartPaymentCalled {
		t.Fatal("expected validation to run")
	}

	if flow.startPaymentCalled {
		t.Fatal("expected StartPayment not to be called")
	}
}
