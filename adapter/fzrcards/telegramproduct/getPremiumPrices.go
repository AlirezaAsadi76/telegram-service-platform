package telegramproduct

import (
	"context"
	"encoding/json"
	"net/http"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"

	"github.com/shopspring/decimal"
)

func (a *Adapter) GetPremiumPlans(
	ctx context.Context,
) ([]productentity.PremiumPrice, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		getPremiumMethod,
		a.createURL(getPremiumPath),
		nil,
	)
	req.Header.Add("x-api-key", a.client.APIKey())

	if err != nil {
		return nil, err
	}

	resp, err := a.client.Connection().Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result premiumResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	prices := make([]productentity.PremiumPrice, 0, len(result.Plans))

	for _, plan := range result.Plans {

		price, pErr := decimal.NewFromString(plan.PriceUSD)

		if pErr != nil {
			return nil, pErr
		}

		prices = append(
			prices,
			productentity.PremiumPrice{
				Months:   uint8(plan.Months),
				PriceUSD: entity.Amount(price),
			},
		)
	}

	return prices, nil
}
