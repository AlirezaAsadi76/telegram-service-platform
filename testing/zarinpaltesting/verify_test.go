package zarinpaltesting

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"telegram-service-platform/adapter/payment/zarinpal"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/pkg/richerror"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

const verifyPath = "/pg/v4/payment/verify.json"

func TestAdapter_Verify_Success(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf(
					"expected POST, got %s",
					r.Method,
				)
			}

			if r.URL.Path != verifyPath {
				t.Fatalf(
					"expected path %s, got %s",
					verifyPath,
					r.URL.Path,
				)
			}

			var request zarinpal.VerifyRequest

			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Fatalf(
					"decode request: %v",
					err,
				)
			}

			if request.MerchantID != "test-merchant" {
				t.Fatalf("unexpected merchant ID")
			}

			if request.Authority != "AUTH-123" {
				t.Fatalf("unexpected authority")
			}

			if request.Amount != 100000 {
				t.Fatalf(
					"expected amount 100000, got %d",
					request.Amount,
				)
			}

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			_ = json.NewEncoder(w).Encode(
				zarinpal.VerifyResponse{
					Data: zarinpal.VerifyData{
						Code:  100,
						RefID: 987654,
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

	result, err := adapter.Verify(
		context.Background(),
		paymentproviderparams.VerifyRequest{
			PaymentID:  100,
			ExternalID: "AUTH-123",
			Amount:     entity.Amount(decimal.NewFromInt(100000)),
			Currency:   entity.CurrencyTOMAN,
		},
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if result.Status != paymententity.PaymentStatusSuccess {
		t.Fatalf(
			"expected success, got %s",
			result.Status,
		)
	}

	if result.ReferenceID != "987654" {
		t.Fatalf(
			"expected reference ID 987654, got %s",
			result.ReferenceID,
		)
	}
}

func TestAdapter_Verify_AlreadyVerified(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			_ = json.NewEncoder(w).Encode(
				zarinpal.VerifyResponse{
					Data: zarinpal.VerifyData{
						Code:  101,
						RefID: 123456,
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

	result, err := adapter.Verify(
		context.Background(),
		paymentproviderparams.VerifyRequest{
			PaymentID:  100,
			ExternalID: "AUTH-123",
			Amount: entity.Amount(
				decimal.NewFromInt(100000),
			),
			Currency: entity.CurrencyTOMAN,
		},
	)

	if err != nil {
		t.Fatalf(
			"expected nil error, got %v",
			err,
		)
	}

	if result.Status != paymententity.PaymentStatusSuccess {
		t.Fatalf(
			"expected SUCCESS, got %s",
			result.Status,
		)
	}

	if result.ReferenceID != "123456" {
		t.Fatalf(
			"unexpected reference ID: %s",
			result.ReferenceID,
		)
	}
}

func TestAdapter_Verify_Unavailable(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(
				w,
				"service unavailable",
				http.StatusServiceUnavailable,
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

	_, err := adapter.Verify(
		context.Background(),
		paymentproviderparams.VerifyRequest{
			ExternalID: "AUTH-123",
			Amount: entity.Amount(
				decimal.NewFromInt(100000),
			),
			Currency: entity.CurrencyTOMAN,
		},
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !richerror.IsCode(err, richerror.CodePaymentProviderUnavailable) {
		t.Fatalf("expected unavailable code, got %v", err)
	}
}

func TestAdapter_Verify_InvalidJSON(t *testing.T) {
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

	_, err := adapter.Verify(
		context.Background(),
		paymentproviderparams.VerifyRequest{
			ExternalID: "AUTH-123",
			Amount: entity.Amount(
				decimal.NewFromInt(100000),
			),
			Currency: entity.CurrencyTOMAN,
		},
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !richerror.IsCode(err, richerror.CodePaymentProviderInvalidResponse) {
		t.Fatalf("expected invalid response code, got %v", err)
	}
}

func TestAdapter_Verify_Rejected(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			_ = json.NewEncoder(w).Encode(
				zarinpal.VerifyResponse{
					Data: zarinpal.VerifyData{
						Code:  1000,
						RefID: 0,
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

	_, err := adapter.Verify(
		context.Background(),
		paymentproviderparams.VerifyRequest{
			ExternalID: "AUTH-123",
			Amount: entity.Amount(
				decimal.NewFromInt(100000),
			),
			Currency: entity.CurrencyTOMAN,
		},
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !richerror.IsCode(err, richerror.CodePaymentProviderRejected) {
		t.Fatalf("expected rejected code, got %v", err)
	}
}

func TestAdapter_Verify_Timeout(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
		}))

	defer server.Close()

	adapter := zarinpal.New(
		zarinpal.Config{
			MerchantID: "test-merchant",
			BaseURL:    server.URL,
			Timeout:    20 * time.Millisecond,
		},
		nil,
	)

	_, err := adapter.Verify(
		context.Background(),
		paymentproviderparams.VerifyRequest{
			ExternalID: "AUTH-123",
			Amount: entity.Amount(
				decimal.NewFromInt(100000),
			),
			Currency: entity.CurrencyTOMAN,
		},
	)

	if err == nil {
		t.Fatal("expected timeout error")
	}

	if !richerror.IsCode(err, richerror.CodePaymentProviderTimeout) {
		t.Fatalf("expected timeout code, got %v", err)
	}
}
