package pricingtesting

import (
	"context"
	"errors"
	"telegram-service-platform/entity"
	"telegram-service-platform/service/pricingservice"
	"testing"

	"github.com/shopspring/decimal"
)

func TestService_CalculatePrice(t *testing.T) {
	tests := []struct {
		name          string
		usd           float64
		tonUsdPrice   float64
		usdTomanPrice float64
		tonUsdErr     error
		usdTomanErr   error
		wantErr       bool
		wantTON       float64
		wantToman     float64
		wantUSDT      float64
	}{
		{
			name:          "success",
			usd:           10,
			tonUsdPrice:   2,
			usdTomanPrice: 100000,
			wantTON:       5,
			wantToman:     1000000,
			wantUSDT:      10,
		},
		{
			name:        "invalid ton price",
			usd:         10,
			tonUsdPrice: 0,
			wantErr:     true,
		},
		{
			name:          "invalid toman price",
			usd:           10,
			tonUsdPrice:   2,
			usdTomanPrice: 0,
			wantErr:       true,
		},
		{
			name:      "ton price repository error",
			usd:       10,
			tonUsdErr: errors.New("ton price unavailable"),
			wantErr:   true,
		},
		{
			name:        "toman price repository error",
			usd:         10,
			tonUsdPrice: 2,
			usdTomanErr: errors.New(
				"toman price unavailable",
			),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakePriceRepository{
				tonUsdPrice:   tt.tonUsdPrice,
				usdTomanPrice: tt.usdTomanPrice,
				tonUsdErr:     tt.tonUsdErr,
				usdTomanErr:   tt.usdTomanErr,
			}

			service := pricingservice.New(repository)

			usd := entity.Amount(
				decimal.NewFromFloat(tt.usd),
			)

			price, err := service.CalculatePrice(
				context.Background(),
				usd,
			)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertAmount(t, price.USD, tt.usd)
			assertAmount(t, price.USDT, tt.wantUSDT)
			assertAmount(t, price.TON, tt.wantTON)
			assertAmount(t, price.Toman, tt.wantToman)
		})
	}
}
