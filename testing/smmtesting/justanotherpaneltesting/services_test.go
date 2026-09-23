package justanotherpaneltesting

import (
	"net/http"
	"net/http/httptest"
	"telegram-service-platform/adapter/smm/justanotherpanel"
	"testing"
	"time"
)

func TestAdapter_AllServices_MapsProviderMetadata(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			if err := r.ParseForm(); err != nil {
				t.Fatalf("parse form: %v", err)
			}

			if r.Form.Get("action") != string(justanotherpanel.ActionTypeServices) {
				t.Fatalf(
					"expected action %s, got %s",
					justanotherpanel.ActionTypeServices,
					r.Form.Get("action"),
				)
			}

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			_, _ = w.Write([]byte(`[
				{
					"service": 1234,
					"name": "Telegram Members",
					"type": "Default",
					"rate": "0.125",
					"min": 10,
					"max": 100000,
					"dripfeed": true,
					"refill": false,
					"cancel": true,
					"category": "Telegram Members"
				}
			]`))
		}),
	)
	defer server.Close()

	adapter := justanotherpanel.New(justanotherpanel.Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})

	response, err := adapter.AllServices(
		t.Context(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(response.Services) != 1 {
		t.Fatalf(
			"expected 1 service, got %d",
			len(response.Services),
		)
	}

	service := response.Services[0]

	if service.Service != 1234 {
		t.Fatalf(
			"expected service ID 1234, got %d",
			service.Service,
		)
	}

	if service.Name != "Telegram Members" {
		t.Fatalf(
			"expected service name Telegram Members, got %s",
			service.Name,
		)
	}

	if !service.IsActive {
		t.Fatal("expected provider service to be active")
	}

	if service.ProviderName != justanotherpanel.ProviderName {
		t.Fatalf(
			"expected provider %s, got %s",
			justanotherpanel.ProviderName,
			service.ProviderName,
		)
	}
}
