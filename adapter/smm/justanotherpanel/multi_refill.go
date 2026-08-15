package justanotherpanel

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"telegram-service-platform/params/smmprams"
)

// MultiRefill requests refills for multiple orders at once.
// JAP returns a map of order_id -> refill_id on success.
// Note: JAP docs show duplicate keys in example, but actual API returns valid JSON map.
func (a *Adapter) MultiRefill(ctx context.Context, orderIDs []string) (smmprams.MultiRefillResponse, error) {
	form := url.Values{
		"key":    {a.config.APIKey},
		"action": {string(ActionTypeRefill)},
		"orders": {strings.Join(orderIDs, ",")},
	}

	var result map[string]int64
	if err := a.doRequest(ctx, form, &result); err != nil {
		return smmprams.MultiRefillResponse{}, fmt.Errorf("multi-refill request: %w", err)
	}

	return smmprams.MultiRefillResponse{Refills: result}, nil
}
