package justanotherpanel

import (
	"context"
	"fmt"
	"net/url"

	"telegram-service-platform/params/smmparams"
)

// Refill requests a refill for a single order.
// JAP returns {"refill": <refill_id>} on success.
func (a *Adapter) Refill(ctx context.Context, orderID string) (smmparams.RefillResponse, error) {
	form := url.Values{
		"key":    {a.config.APIKey},
		"action": {string(ActionTypeRefill)},
		"order":  {orderID},
	}

	var result struct {
		Refill int64  `json:"refill"`
		Error  string `json:"error"`
	}
	if err := a.doRequest(ctx, form, &result); err != nil {
		return smmparams.RefillResponse{}, fmt.Errorf("refill request: %w", err)
	}

	if result.Error != "" {
		return smmparams.RefillResponse{}, fmt.Errorf("jap error: %s", result.Error)
	}

	return smmparams.RefillResponse{RefillID: result.Refill}, nil
}
