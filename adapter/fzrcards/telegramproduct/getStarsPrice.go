package telegramproduct

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

func (a *Adapter) GetStarPrice(ctx context.Context) (float64, error) {
	const Op = "telegramproduct.GetStarPrice"

	req, nErr := http.NewRequestWithContext(ctx, getStarsPriceMethod, a.createURL(getStarsPricePath), nil)

	if nErr != nil {
		return 0, nErr
	}

	resp, cErr := a.client.Connection().Do(req)
	if cErr != nil {
		return 0, cErr
	}

	defer resp.Body.Close()

	var result starsPriceResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(
		result.PricePerStar,
		64,
	)

	if err != nil {
		return 0, err
	}

	return price, nil
}
