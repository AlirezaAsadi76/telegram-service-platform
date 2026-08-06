package exchangerate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (e *Adapter) GetUsdTomanPrice(ctx context.Context) (float64, error) {
	const Op = "exchangerate.getUsdTomanPrice"

	req, nErr := http.NewRequestWithContext(ctx, http.MethodGet, e.config.UsdIrURL, nil)
	req.Header.Set("Accept", "application/json")
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

	body, eErr := io.ReadAll(resp.Body)

	if eErr != nil {
		return 0, eErr
	}

	var UsdTomanPrice marketsResponse

	deErr := json.NewDecoder(bytes.NewReader(body)).Decode(&UsdTomanPrice)

	if deErr != nil {
		return 0, deErr
	}
	usdtmn, ok := UsdTomanPrice.Result[usdTomanMarket]

	if !ok {
		return 0, errors.New("USDTTMN market not found")
	}

	return float64(usdtmn.Stats.LastPrice), nil
}
