package producttesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/params/productparams"
	"telegram-service-platform/params/smmparams"
	"telegram-service-platform/service/productservice"
	"testing"

	"github.com/shopspring/decimal"
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

func TestService_CalculateSMMPrice(t *testing.T) {
	repository := &calculateSMMPriceRepository{
		mapping: &smmentity.SmmMapping{
			Id:           55,
			SmmServiceId: 10,
			IsActive:     true,
		},
		service: &smmentity.SMM{
			Id:       10,
			Service:  123456,
			Rate:     entity.Amount(decimal.NewFromFloat(0.125)),
			IsActive: true,
		},
	}

	pricingService := newSMMPricingService(
		2,
		100000,
	)

	service := productservice.New(
		productservice.Config{},
		pricingService,
		repository,
		nil,
		nil,
		nil,
	)

	response, err := service.CalculateSMMPrice(
		context.Background(),
		productparams.CalculateSMMPriceRequest{
			MappingID: 55,
			Quantity:  2000,
		},
	)
	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if response.MappingID != 55 {
		t.Fatalf(
			"expected mapping ID 55, got %d",
			response.MappingID,
		)
	}

	if response.ServiceID != 10 {
		t.Fatalf(
			"expected internal service ID 10, got %d",
			response.ServiceID,
		)
	}

	if !response.Rate.Equal(
		entity.Amount(decimal.NewFromFloat(0.125)),
	) {
		t.Fatalf(
			"expected rate 0.125, got %s",
			response.Rate.String(),
		)
	}

	// 0.125 * 2000 / 1000 = 0.25 USD
	expectedUSD := entity.Amount(
		decimal.NewFromFloat(0.25),
	)

	if !response.Price.USD.Equal(expectedUSD) {
		t.Fatalf(
			"expected USD price %s, got %s",
			expectedUSD.String(),
			response.Price.USD.String(),
		)
	}

	// 0.25 / 2 = 0.125 TON
	expectedTON := entity.Amount(
		decimal.NewFromFloat(0.125),
	)

	if !response.Price.TON.Equal(expectedTON) {
		t.Fatalf(
			"expected TON price %s, got %s",
			expectedTON.String(),
			response.Price.TON.String(),
		)
	}

	// 0.25 * 100000 = 25000 TOMAN
	expectedToman := entity.Amount(
		decimal.NewFromFloat(25000),
	)

	if !response.Price.Toman.Equal(expectedToman) {
		t.Fatalf(
			"expected TOMAN price %s, got %s",
			expectedToman.String(),
			response.Price.Toman.String(),
		)
	}
}

func TestService_CalculateSMMPrice_InactiveMapping(
	t *testing.T,
) {
	repository := &calculateSMMPriceRepository{
		mapping: &smmentity.SmmMapping{
			Id:           55,
			SmmServiceId: 10,
			IsActive:     false,
		},
		service: &smmentity.SMM{
			Id:       10,
			Rate:     entity.Amount(decimal.NewFromFloat(0.125)),
			IsActive: true,
		},
	}

	service := productservice.New(
		productservice.Config{},
		newSMMPricingService(2, 100000),
		repository,
		nil,
		nil,
		nil,
	)

	_, err := service.CalculateSMMPrice(
		context.Background(),
		productparams.CalculateSMMPriceRequest{
			MappingID: 55,
			Quantity:  1000,
		},
	)
	if err == nil {
		t.Fatal("expected inactive mapping to return error")
	}
}

func TestService_CalculateSMMPrice_InactiveService(
	t *testing.T,
) {
	repository := &calculateSMMPriceRepository{
		mapping: &smmentity.SmmMapping{
			Id:           55,
			SmmServiceId: 10,
			IsActive:     true,
		},
		service: &smmentity.SMM{
			Id:       10,
			Rate:     entity.Amount(decimal.NewFromFloat(0.125)),
			IsActive: false,
		},
	}

	service := productservice.New(
		productservice.Config{},
		newSMMPricingService(2, 100000),
		repository,
		nil,
		nil,
		nil,
	)

	_, err := service.CalculateSMMPrice(
		context.Background(),
		productparams.CalculateSMMPriceRequest{
			MappingID: 55,
			Quantity:  1000,
		},
	)
	if err == nil {
		t.Fatal("expected inactive SMM service to return error")
	}
}

func TestService_GetSMMServiceByMappingID_ResolvesProviderService(
	t *testing.T,
) {
	repository := &mappingResolutionRepository{
		mapping: &smmentity.SmmMapping{
			Id:           55,
			SmmServiceId: 10,
			IsActive:     true,
		},
		service: &smmentity.SMM{
			Id:           10,
			Service:      123456,
			Name:         "Telegram Members",
			ProviderName: "justanotherpanel",
			Rate: entity.Amount(
				decimal.NewFromFloat(0.125),
			),
			IsActive: true,
		},
	}

	service := productservice.New(
		productservice.Config{},
		nil,
		repository,
		nil,
		nil,
		nil,
	)

	smm, err := service.GetSMMServiceByMappingID(
		context.Background(),
		55,
	)
	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if smm.Id != 10 {
		t.Fatalf(
			"expected internal SMM ID 10, got %d",
			smm.Id,
		)
	}

	if smm.Service != 123456 {
		t.Fatalf(
			"expected provider service ID 123456, got %d",
			smm.Service,
		)
	}

	if smm.ProviderName != "justanotherpanel" {
		t.Fatalf(
			"expected provider justanotherpanel, got %s",
			smm.ProviderName,
		)
	}
}

func TestService_GetSMMServiceByMappingID_MappingNotFound(
	t *testing.T,
) {
	repository := &mappingResolutionRepository{
		mappingErr: errors.New("mapping not found"),
	}

	service := productservice.New(
		productservice.Config{},
		nil,
		repository,
		nil,
		nil,
		nil,
	)

	_, err := service.GetSMMServiceByMappingID(
		context.Background(),
		55,
	)
	if err == nil {
		t.Fatal("expected mapping lookup error")
	}
}

func TestService_GetSMMServiceByMappingID_ServiceNotFound(
	t *testing.T,
) {
	repository := &mappingResolutionRepository{
		mapping: &smmentity.SmmMapping{
			Id:           55,
			SmmServiceId: 10,
		},
		serviceErr: errors.New(
			"SMM service not found",
		),
	}

	service := productservice.New(
		productservice.Config{},
		nil,
		repository,
		nil,
		nil,
		nil,
	)

	_, err := service.GetSMMServiceByMappingID(
		context.Background(),
		55,
	)
	if err == nil {
		t.Fatal("expected service lookup error")
	}
}
