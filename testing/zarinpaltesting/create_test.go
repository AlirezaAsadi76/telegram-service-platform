package zarinpaltesting

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"telegram-service-platform/adapter/payment/zarinpal"
	"telegram-service-platform/entity"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/pkg/richerror"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

const createPath = "/pg/v4/payment/request.json"

func TestAdapter_Create_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != createPath {
			t.Fatalf("expected %s, got %s", createPath, r.URL.Path)
		}

		var request zarinpal.CreateRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		if request.MerchantID != "test-merchant" {
			t.Fatalf("unexpected merchant id")
		}

		if request.Amount != 100000 {
			t.Fatalf("expected amount 100000, got %d", request.Amount)
		}

		if request.CallbackURL != "https://example.com/callback" {
			t.Fatalf("unexpected callback URL")
		}

		if request.Description != "test payment" {
			t.Fatalf("unexpected description")
		}

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		_ = json.NewEncoder(w).Encode(
			zarinpal.CreateResponse{
				Data: zarinpal.CreateData{
					Code:      100,
					Authority: "AUTH-123",
				},
			},
		)
	}),
	)

	defer server.Close()

	adapter := zarinpal.New(
		zarinpal.Config{
			MerchantID:  "test-merchant",
			BaseURL:     server.URL,
			StartPayURL: server.URL + "/StartPay",
		},
		server.Client(),
	)

	result, err := adapter.Create(
		context.Background(),
		paymentproviderparams.CreateRequest{
			PaymentID:   100,
			OrderID:     10,
			Amount:      entity.Amount(decimal.NewFromInt(100000)),
			Currency:    entity.CurrencyTOMAN,
			CallbackURL: "https://example.com/callback",
			Description: "test payment",
		},
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if result.ExternalID != "AUTH-123" {
		t.Fatalf(
			"expected external ID AUTH-123, got %s",
			result.ExternalID,
		)
	}

	expectedURL := server.URL + "/StartPay/AUTH-123"

	if result.PaymentURL != expectedURL {
		t.Fatalf("expected payment URL %s, got %s", expectedURL, result.PaymentURL)
	}
}

func TestAdapter_Create_Rejected(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			_ = json.NewEncoder(w).Encode(
				zarinpal.CreateResponse{
					Data: zarinpal.CreateData{
						Code: 0,
					},
					Errors: &zarinpal.CreateErrors{
						Code:    1001,
						Message: "invalid merchant",
					},
				},
			)
		}),
	)

	defer server.Close()

	adapter := zarinpal.New(
		zarinpal.Config{
			MerchantID: "test-merchant",
			BaseURL:    server.URL,
		},
		server.Client(),
	)

	_, err := adapter.Create(
		context.Background(),
		paymentproviderparams.CreateRequest{
			Amount:      entity.Amount(decimal.NewFromInt(100000)),
			Currency:    entity.CurrencyTOMAN,
			CallbackURL: "https://example.com/callback",
		},
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !richerror.IsCode(err, richerror.CodePaymentProviderRejected) {
		t.Fatalf("expected rejected code, got %v", err)
	}
}

func TestAdapter_Create_Timeout(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
		}),
	)

	defer server.Close()

	adapter := zarinpal.New(
		zarinpal.Config{
			MerchantID: "test-merchant",
			BaseURL:    server.URL,
			Timeout:    20 * time.Millisecond,
		},
		nil,
	)

	_, err := adapter.Create(
		context.Background(),
		paymentproviderparams.CreateRequest{
			Amount:      entity.Amount(decimal.NewFromInt(100000)),
			Currency:    entity.CurrencyTOMAN,
			CallbackURL: "https://example.com/callback",
		},
	)

	if err == nil {
		t.Fatal("expected timeout error")
	}

	if !richerror.IsCode(err, richerror.CodePaymentProviderTimeout) {
		t.Fatalf("expected timeout code, got %v", err)
	}
}

func TestAdapter_Create_InvalidResponse(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			_, _ = w.Write(
				[]byte(`{"data":`),
			)
		}),
	)

	defer server.Close()

	adapter := zarinpal.New(
		zarinpal.Config{
			MerchantID: "test-merchant",
			BaseURL:    server.URL,
		},
		server.Client(),
	)

	_, err := adapter.Create(
		context.Background(),
		paymentproviderparams.CreateRequest{
			Amount:      entity.Amount(decimal.NewFromInt(100000)),
			Currency:    entity.CurrencyTOMAN,
			CallbackURL: "https://example.com/callback",
		},
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !richerror.IsCode(err, richerror.CodePaymentProviderInvalidResponse) {
		t.Fatalf("expected invalid response code, got %v", err)
	}
}

func TestAdapter_Create_EmptyAuthority(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			_ = json.NewEncoder(w).Encode(
				zarinpal.CreateResponse{
					Data: zarinpal.CreateData{
						Code:      100,
						Authority: "",
					},
				},
			)
		}),
	)

	defer server.Close()

	adapter := zarinpal.New(
		zarinpal.Config{
			MerchantID: "test-merchant",
			BaseURL:    server.URL,
		},
		server.Client(),
	)

	_, err := adapter.Create(
		context.Background(),
		paymentproviderparams.CreateRequest{
			Amount:      entity.Amount(decimal.NewFromInt(100000)),
			Currency:    entity.CurrencyTOMAN,
			CallbackURL: "https://example.com/callback",
		},
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !richerror.IsCode(err, richerror.CodePaymentProviderInvalidResponse) {
		t.Fatalf("expected invalid response code, got %v", err)
	}
}

func TestAdapter_Create_UnsupportedCurrency(t *testing.T) {
	called := false

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		}),
	)

	defer server.Close()

	adapter := zarinpal.New(
		zarinpal.Config{
			MerchantID: "test-merchant",
			BaseURL:    server.URL,
		},
		server.Client(),
	)

	_, err := adapter.Create(
		context.Background(),
		paymentproviderparams.CreateRequest{
			Amount:   entity.Amount(decimal.NewFromInt(100000)),
			Currency: entity.CurrencyTON,
		},
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if !richerror.IsCode(err, richerror.CodePaymentUnsupportedCurrency) {
		t.Fatalf("unexpected code: %v", err)
	}

	if called {
		t.Fatal("HTTP request should not be sent")
	}
}
