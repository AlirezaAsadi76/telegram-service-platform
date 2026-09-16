package justanotherpanel

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/params/smmparams"
)

// MultiStatus queries the status of multiple orders at once.
// JAP returns a map where keys are order IDs and values are status objects.
func (a *Adapter) MultiStatus(ctx context.Context, orderIDs []string) (smmparams.GetMultiStatusResponse, error) {
	form := url.Values{
		"key":    {a.config.APIKey},
		"action": {string(ActionTypeStatus)},
		"orders": {strings.Join(orderIDs, ",")},
	}

	// JAP returns: {"1234": {...}, "5678": {...}}
	var result map[string]smmentity.Status
	if err := a.doRequest(ctx, form, &result); err != nil {
		return smmparams.GetMultiStatusResponse{}, fmt.Errorf("multi-status request: %w", err)
	}

	return smmparams.GetMultiStatusResponse{Statuses: result}, nil
}
