package telegramproduct

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"telegram-service-platform/entity/productentity"
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

	price, err := strconv.ParseFloat(
		result.PricePerStar,
		64,
	)

	if err != nil {
		return productentity.StarPrice{}, err
	}

	return productentity.StarPrice{
		PricePerStar: price,
	}, nil
}
