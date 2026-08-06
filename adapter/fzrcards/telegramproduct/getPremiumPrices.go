package telegramproduct

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"telegram-service-platform/pkg/productprice"
)

func (a Adapter) GetPremiumPrices(
	ctx context.Context,
) ([]productprice.PremiumPrice, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		getPremiumMethod,
		a.createURL(getPremiumPath),
		nil,
	)

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

	prices := make([]productprice.PremiumPrice, 0, len(result.Plans))

	for _, plan := range result.Plans {

		price, err := strconv.ParseFloat(
			plan.PriceUSD,
			64,
		)

		if err != nil {
			return nil, err
		}

		prices = append(
			prices,
			productprice.PremiumPrice{
				Months:   plan.Months,
				PriceUSD: price,
			},
		)
	}

	return prices, nil
}
