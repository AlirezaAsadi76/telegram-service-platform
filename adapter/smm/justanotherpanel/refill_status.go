package justanotherpanel

import (
	"context"
	"fmt"
	"net/url"

	"telegram-service-platform/params/smmparams"
)

// RefillStatus queries the status of a single refill request.
// JAP returns {"status": "Completed"} or {"error": "Refill not found"}.
func (a *Adapter) RefillStatus(ctx context.Context, refillID string) (smmparams.RefillStatusResponse, error) {
	form := url.Values{
		"key":    {a.config.APIKey},
		"action": {string(ActionTypeRefillStatus)},
		"refill": {refillID},
	}

	var result struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := a.doRequest(ctx, form, &result); err != nil {
		return smmparams.RefillStatusResponse{}, fmt.Errorf("refill-status request: %w", err)
	}

	if result.Error != "" {
		return smmparams.RefillStatusResponse{Error: result.Error}, nil
	}

	return smmparams.RefillStatusResponse{Status: result.Status}, nil
}
