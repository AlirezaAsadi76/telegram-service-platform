package justanotherpanel

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/params/smmprams"
)

// MultiStatus queries the status of multiple orders at once.
// JAP returns a map where keys are order IDs and values are status objects.
func (a *Adapter) MultiStatus(ctx context.Context, orderIDs []string) (smmprams.GetMultiStatusResponse, error) {
	form := url.Values{
		"key":    {a.config.APIKey},
		"action": {string(ActionTypeStatus)},
		"orders": {strings.Join(orderIDs, ",")},
	}

	// JAP returns: {"1234": {...}, "5678": {...}}
	var result map[string]smmentity.Status
	if err := a.doRequest(ctx, form, &result); err != nil {
		return smmprams.GetMultiStatusResponse{}, fmt.Errorf("multi-status request: %w", err)
	}

	return smmprams.GetMultiStatusResponse{Statuses: result}, nil
}
