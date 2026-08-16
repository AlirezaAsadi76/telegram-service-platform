package justanotherpanel

import (
	"context"
	"fmt"
	"net/url"

	"telegram-service-platform/params/smmprams"
)

// RefillStatus queries the status of a single refill request.
// JAP returns {"status": "Completed"} or {"error": "Refill not found"}.
func (a *Adapter) RefillStatus(ctx context.Context, refillID string) (smmprams.RefillStatusResponse, error) {
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
		return smmprams.RefillStatusResponse{}, fmt.Errorf("refill-status request: %w", err)
	}

	if result.Error != "" {
		return smmprams.RefillStatusResponse{Error: result.Error}, nil
	}

	return smmprams.RefillStatusResponse{Status: result.Status}, nil
}
