package exchangerate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (e *Adapter) GetUsdTomanPrice(ctx context.Context) (float64, error) {
	const Op = "exchangerate.getUsdTomanPrice"

	req, nErr := http.NewRequestWithContext(ctx, http.MethodGet, wallexUsdIRREndPoint, nil)
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

	var UsdTomanPrice marketsResponse

	deErr := json.NewDecoder(resp.Body).Decode(&UsdTomanPrice)
	if deErr != nil {
		return 0, deErr
	}
	usdtmn, ok := UsdTomanPrice.Result[usdTomanMarket]

	if !ok {
		return 0, errors.New("USDTTMN market not found")
	}

	return usdtmn.Stats.LastPrice, nil
}
