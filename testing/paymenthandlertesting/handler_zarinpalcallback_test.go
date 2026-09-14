package paymenthandlertesting

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"telegram-service-platform/delivery/httpserver/paymenthandler"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentparams"
	"telegram-service-platform/pkg/richerror"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestHandler_ZarinpalCallback_Success(t *testing.T) {
	paymentService := &fakePaymentFlow{
		response: &paymentparams.ConfirmPaymentResponse{
			PaymentID: 100,
			OrderID:   10,
			Status:    paymententity.PaymentStatusSuccess,
		},
	}

	validator := &fakePaymentValidator{}

	handler := paymenthandler.New(
		paymentService,
		validator,
		MiddlewareTest{},
	)

	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/payments/zarinpal/callback?status=OK&authority=AUTH-123",
		nil,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := handler.ZarinpalCallbackHandler(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if validator.calls != 1 {
		t.Fatalf("expected validator to be called once, got %d", validator.calls)
	}

	if paymentService.calls != 1 {
		t.Fatalf("expected service to be called once, got %d", paymentService.calls)
	}

	if paymentService.lastReq.ExternalID != "AUTH-123" {
		t.Fatalf("expected AUTH-123, got %s", paymentService.lastReq.ExternalID)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_ZarinpalCallback_MissingAuthority(t *testing.T) {
	paymentService := &fakePaymentFlow{}

	validator := &fakePaymentValidator{
		err: richerror.New(
			"paymentvalidator.zarinpalcallback",
			nil,
		).
			WithKind(richerror.KindInvalid),
		fieldErrs: map[string]string{
			"authority": "is required",
		},
	}

	handler := paymenthandler.New(
		paymentService,
		validator,
		MiddlewareTest{},
	)

	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/payments/zarinpal/callback?status=OK",
		nil,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := handler.ZarinpalCallbackHandler(c)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}

	if paymentService.calls != 0 {
		t.Fatalf("expected service not to be called, got %d", paymentService.calls)
	}
}

func TestHandler_ZarinpalCallback_NOK(t *testing.T) {
	paymentService := &fakePaymentFlow{}
	validator := &fakePaymentValidator{}

	handler := paymenthandler.New(
		paymentService,
		validator,
		MiddlewareTest{},
	)

	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/payments/zarinpal/callback?status=NOK&authority=AUTH-123",
		nil,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := handler.ZarinpalCallbackHandler(c)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if paymentService.calls != 0 {
		t.Fatalf("expected service not to be called, got %d", paymentService.calls)
	}
}

func TestHandler_ZarinpalCallback_ServiceError(
	t *testing.T,
) {
	paymentService := &fakePaymentFlow{
		err: richerror.New(
			"paymentservice.confirm",
			errors.New("payment not found"),
		).
			WithKind(richerror.KindNotFound).
			WithCode(richerror.CodePaymentNotFound),
	}

	validator := &fakePaymentValidator{}

	handler := paymenthandler.New(
		paymentService,
		validator,
		MiddlewareTest{},
	)

	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/payments/zarinpal/callback?status=OK&authority=AUTH-123",
		nil,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := handler.ZarinpalCallbackHandler(c)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !richerror.IsCode(err, richerror.CodePaymentNotFound) {
		t.Fatalf("expected payment not found, got %v", err)
	}
}
func TestHandler_SetRoutes(t *testing.T) {
	handler := paymenthandler.New(
		&fakePaymentFlow{},
		&fakePaymentValidator{},
		MiddlewareTest{},
	)

	e := echo.New()

	handler.SetRoutes(e)

	routes := e.Router().Routes()

	found := false

	for _, route := range routes {
		if route.Path == "/payments/zarinpal/callback" &&
			route.Method == http.MethodGet {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("zarinpal callback route was not registered")
	}
}
