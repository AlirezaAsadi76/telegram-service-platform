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

			service := pricingservice.New(repository, pricingservice.Config{
				SMMMinimumUnitPriceToman: "5000",
				SMMMinimumMultiplier:     "1.4",
				SMMMaximumMultiplier:     "4",
				SMMMultiplierScaleToman:  "1466.6666666666666666666666666667",
			})

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

func TestService_CalculateSMMSellingPrice_UsesMinimumPriceForCheapService(
	t *testing.T,
) {
	service := newSMMPricingServiceWithConfig(
		pricingservice.Config{
			SMMMinimumUnitPriceToman: "5000",
			SMMMinimumMultiplier:     "1.4",
			SMMMaximumMultiplier:     "4",
			SMMMultiplierScaleToman:  "1466.6666666666666666666666666667",
		},
		2,
		240000,
	)

	price, err := service.CalculateSMMSellingPrice(
		context.Background(),
		entity.Amount(decimal.NewFromFloat(0.0015)),
		15000,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Provider cost:
	// 0.0015 * 15000 / 1000 = 0.0225 USD
	//
	// Base cost:
	// 0.0225 * 240000 = 5400 TOMAN
	//
	// Base unit cost:
	// 5400 / 15 = 360 TOMAN / 1000
	//
	// Minimum selling price:
	// 5000 * 15 = 75000 TOMAN
	expected := entity.Amount(
		decimal.NewFromInt(75000),
	)

	if !price.Toman.Equal(expected) {
		t.Fatalf(
			"expected TOMAN price %s, got %s",
			expected.String(),
			price.Toman.String(),
		)
	}
}

func TestService_CalculateSMMSellingPrice_IsContinuousAtMinimumPriceBoundary(
	t *testing.T,
) {
	service := newSMMPricingServiceWithConfig(
		pricingservice.Config{
			SMMMinimumUnitPriceToman: "5000",
			SMMMinimumMultiplier:     "1.4",
			SMMMaximumMultiplier:     "4",
			SMMMultiplierScaleToman:  "1466.6666666666666666666666666667",
		},
		2,
		240000,
	)

	price, err := service.CalculateSMMSellingPrice(
		context.Background(),
		entity.Amount(decimal.NewFromInt(2000)).Div(
			entity.Amount(decimal.NewFromInt(240000)),
		),
		1000,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := entity.Amount(
		decimal.NewFromInt(5000),
	)

	if !price.Toman.Equal(expected) {
		t.Fatalf(
			"expected boundary price %s, got %s",
			expected.String(),
			price.Toman.String(),
		)
	}
}

func TestSMMPricingRule_CalculatesExpectedSellingUnitPrice(
	t *testing.T,
) {
	rule := pricingservice.SMMPricingRule{
		MinimumUnitPriceToman: entity.Amount(
			decimal.NewFromInt(5000),
		),
		MinimumMultiplier: decimal.RequireFromString("1.4"),
		MaximumMultiplier: decimal.RequireFromString("4"),
		MultiplierScaleToman: entity.Amount(
			decimal.RequireFromString(
				"1466.6666666666666666666666666667",
			),
		),
	}

	price, err := rule.CalculateSellingUnitPrice(
		entity.Amount(decimal.NewFromInt(10000)),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := entity.Amount(
		decimal.NewFromInt(17326),
	)

	if !price.Equal(expected) {
		t.Fatalf(
			"expected unit price %s, got %s",
			expected.String(),
			price.String(),
		)
	}
}
