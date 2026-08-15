package justanotherpanel

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"telegram-service-platform/params/smmprams"
)

// MultiRefillStatus queries the status of multiple refill requests at once.
// JAP returns: {"1": {"status": "Completed"}, "2": {"status": "Rejected"}}
func (a *Adapter) MultiRefillStatus(ctx context.Context, refillIDs []string) (smmprams.MultiRefillStatusResponse, error) {
	form := url.Values{
		"key":     {a.config.APIKey},
		"action":  {string(ActionTypeRefillStatus)},
		"refills": {strings.Join(refillIDs, ",")},
	}

	var result map[string]struct {
		Status string `json:"status"`
	}
	if err := a.doRequest(ctx, form, &result); err != nil {
		return smmprams.MultiRefillStatusResponse{}, fmt.Errorf("multi-refill-status request: %w", err)
	}

	statuses := make(map[string]string, len(result))
	for id, v := range result {
		statuses[id] = v.Status
	}

	return smmprams.MultiRefillStatusResponse{Statuses: statuses}, nil
}
