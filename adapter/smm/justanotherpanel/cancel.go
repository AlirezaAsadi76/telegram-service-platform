package justanotherpanel

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"telegram-service-platform/params/smmprams"
)

// Cancel requests cancellation for one or more orders.
// JAP returns a map of order_id -> cancel_status (1 = success, 0 = failed).
func (a *Adapter) Cancel(ctx context.Context, orderIDs []string) (smmprams.CancelResponse, error) {
	form := url.Values{
		"key":    {a.config.APIKey},
		"action": {string(ActionTypeCancel)},
		"orders": {strings.Join(orderIDs, ",")},
	}

	var result map[string]int64
	if err := a.doRequest(ctx, form, &result); err != nil {
		return smmprams.CancelResponse{}, fmt.Errorf("cancel request: %w", err)
	}

	return smmprams.CancelResponse{Cancels: result}, nil
}
