package producttesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/params/smmparams"
	"testing"
)

func TestService_SyncSMMServices(t *testing.T) {
	tests := []struct {
		name              string
		response          smmparams.GetAllServicesResponse
		adapterErr        error
		failOnServiceID   int64
		wantErr           bool
		wantUpsertedCount int
	}{
		{
			name: "provider success",
			response: smmparams.GetAllServicesResponse{
				Services: []smmentity.SMM{
					{
						Service:      100,
						Name:         "Telegram Members",
						ProviderName: "justanotherpanel",
						IsActive:     true,
					},
					{
						Service:      200,
						Name:         "Telegram Views",
						ProviderName: "justanotherpanel",
						IsActive:     true,
					},
				},
			},
			wantUpsertedCount: 2,
		},
		{
			name:       "provider error",
			adapterErr: errors.New("provider unavailable"),
			wantErr:    true,
		},
		{
			name: "empty provider catalog",
			response: smmparams.GetAllServicesResponse{
				Services: []smmentity.SMM{},
			},
			wantUpsertedCount: 0,
		},
		{
			name: "repository error",
			response: smmparams.GetAllServicesResponse{
				Services: []smmentity.SMM{
					{
						Service:      100,
						Name:         "Telegram Members",
						ProviderName: "justanotherpanel",
						IsActive:     true,
					},
					{
						Service:      200,
						Name:         "Telegram Views",
						ProviderName: "justanotherpanel",
						IsActive:     true,
					},
				},
			},
			failOnServiceID:   200,
			wantErr:           true,
			wantUpsertedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakeProductRepository{
				failOnServiceID: tt.failOnServiceID,
			}

			adapter := &fakeSMMAdapter{
				response: tt.response,
				err:      tt.adapterErr,
			}

			service := newProductService(
				repository,
				adapter,
			)

			err := service.SyncSMMServices(
				context.Background(),
			)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf(
					"unexpected error: %v",
					err,
				)
			}

			if adapter.calls != 1 {
				t.Fatalf(
					"expected provider to be called once, got %d",
					adapter.calls,
				)
			}

			if len(repository.upserted) != tt.wantUpsertedCount {
				t.Fatalf(
					"expected %d upserts, got %d",
					tt.wantUpsertedCount,
					len(repository.upserted),
				)
			}
		})
	}
}

func TestService_GetMissingSMMServices_ReturnsOnlyMissingServices(
	t *testing.T,
) {
	repository := &fakeProductRepository{
		services: []smmentity.SMM{
			{
				Service:      100,
				Name:         "Telegram Members",
				ProviderName: "justanotherpanel",
				IsActive:     true,
			},
			{
				Service:      200,
				Name:         "Telegram Views",
				ProviderName: "justanotherpanel",
				IsActive:     true,
			},
			{
				Service:      300,
				Name:         "Instagram Followers",
				ProviderName: "justanotherpanel",
				IsActive:     true,
			},
		},
	}

	adapter := &fakeSMMAdapter{
		response: smmparams.GetAllServicesResponse{
			Services: []smmentity.SMM{
				{
					Service: 100,
					Name:    "Telegram Members",
				},
				{
					Service: 200,
					Name:    "Telegram Views",
				},
			},
		},
	}

	service := newProductService(
		repository,
		adapter,
	)

	missing, err := service.GetMissingSMMServices(
		context.Background(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(missing) != 1 {
		t.Fatalf(
			"expected 1 missing service, got %d",
			len(missing),
		)
	}

	if missing[0].Service != 300 {
		t.Fatalf(
			"expected missing service 300, got %d",
			missing[0].Service,
		)
	}

	if repository.getAllCalls != 1 {
		t.Fatalf(
			"expected database catalog to be read once, got %d",
			repository.getAllCalls,
		)
	}

	if adapter.calls != 1 {
		t.Fatalf(
			"expected provider catalog to be read once, got %d",
			adapter.calls,
		)
	}
}
