package justanotherpanel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/params/smmprams"
)

// Cancel requests cancellation for one or more orders.
// JAP returns an array: [
//
//	{"order": 9, "cancel": {"error": "Incorrect order ID"}},
//	{"order": 2, "cancel": 1}
//
// ]
// The "cancel" field is int64 (1 = success, 0 = failed) on success,
// or {"error": "..."} on failure for that specific order.
func (a *Adapter) Cancel(ctx context.Context, orderIDs []string) (smmprams.CancelResponse, error) {
	form := url.Values{
		"key":    {a.config.APIKey},
		"action": {string(ActionTypeCancel)},
		"orders": {strings.Join(orderIDs, ",")},
	}

	// JAP returns array where "cancel" can be int64 or {"error": "..."}
	var rawItems []struct {
		Order  int64           `json:"order"`
		Cancel json.RawMessage `json:"cancel"`
	}
	if err := a.doRequest(ctx, form, &rawItems); err != nil {
		return smmprams.CancelResponse{}, fmt.Errorf("cancel request: %w", err)
	}

	items := make([]smmentity.CancelItem, 0, len(rawItems))
	for _, r := range rawItems {
		item := smmentity.CancelItem{Order: r.Order}

		// Try parse as int64 first (success case: 1 = cancelled, 0 = failed)
		var cancelStatus int64
		if err := json.Unmarshal(r.Cancel, &cancelStatus); err == nil {
			item.Cancelled = cancelStatus
		} else {
			// Parse as error object: {"error": "..."}
			var errObj struct {
				Error string `json:"error"`
			}
			if err := json.Unmarshal(r.Cancel, &errObj); err == nil {
				item.Error = errObj.Error
			} else {
				item.Error = "unknown cancel response"
			}
		}
		items = append(items, item)
	}

	return smmprams.CancelResponse{Items: items}, nil
}
