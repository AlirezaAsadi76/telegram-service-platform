package exchangerate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (e *Adapter) GetTonUsdPrice(ctx context.Context) (float64, error) {
	const Op = "exchangerate.GetTonUsdPrice"

	req, nErr := http.NewRequestWithContext(ctx, http.MethodGet, coinGeckoTonUsdEndPoint, nil)

	if nErr != nil {
		return 0, nErr
	}
	resp, dErr := e.client.Do(req)
	if dErr != nil {
		return 0, dErr
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("%s %s returned %s", req.Method, req.URL, resp.Status)
	}

	var tonUSDPrice tonUsdPriceResponse

	deErr := json.NewDecoder(resp.Body).Decode(&tonUSDPrice)
	if deErr != nil {
		return 0, deErr
	}

	return tonUSDPrice[indexTonUsd].USD, nil
}
