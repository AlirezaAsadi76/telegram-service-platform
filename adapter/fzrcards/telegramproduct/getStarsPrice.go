package telegramproduct

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"telegram-service-platform/entity"
	"telegram-service-platform/entity/productentity"

	"github.com/shopspring/decimal"
)

func (a *Adapter) GetStarPrice(ctx context.Context) (productentity.StarPrice, error) {
	const Op = "telegramproduct.GetStarPrice"

	req, nErr := http.NewRequestWithContext(ctx, getStarsPriceMethod, a.createURL(getStarsPricePath), nil)
	req.Header.Add("x-api-key", a.client.APIKey())

	fmt.Println(a.createURL(getStarsPricePath))

	if nErr != nil {
		return productentity.StarPrice{}, nErr
	}

	resp, cErr := a.client.Connection().Do(req)

	if cErr != nil {
		return productentity.StarPrice{}, cErr
	}

	defer resp.Body.Close()

	var result starsPriceResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return productentity.StarPrice{}, err
	}

	price, err := decimal.NewFromString(result.PricePerStar)

	if err != nil {
		return productentity.StarPrice{}, err
	}

	return productentity.StarPrice{
		PricePerStar: entity.Amount(price),
	}, nil
}
