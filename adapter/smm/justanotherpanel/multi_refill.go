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

// MultiRefill requests refills for multiple orders at once.
// JAP returns an array where each item has {"order": <id>, "refill": <refill_id>}
// or {"order": <id>, "refill": {"error": "..."}} on failure for that order.
func (a *Adapter) MultiRefill(ctx context.Context, orderIDs []string) (smmprams.MultiRefillResponse, error) {
	form := url.Values{
		"key":    {a.config.APIKey},
		"action": {string(ActionTypeRefill)},
		"orders": {strings.Join(orderIDs, ",")},
	}

	// JAP returns: [{"order": 1, "refill": 1}, {"order": 2, "refill": {"error": "..."}}]
	var rawItems []struct {
		Order  int64           `json:"order"`
		Refill json.RawMessage `json:"refill"`
	}
	if err := a.doRequest(ctx, form, &rawItems); err != nil {
		return smmprams.MultiRefillResponse{}, fmt.Errorf("multi-refill request: %w", err)
	}

	items := make([]smmentity.MultiRefillItem, 0, len(rawItems))
	for _, r := range rawItems {
		item := smmentity.MultiRefillItem{Order: r.Order}

		// Try parse as int64 first (success case)
		var refillID int64
		if err := json.Unmarshal(r.Refill, &refillID); err == nil {
			item.RefillID = refillID
		} else {
			// Parse as error object: {"error": "..."}
			var errObj struct {
				Error string `json:"error"`
			}
			if err := json.Unmarshal(r.Refill, &errObj); err == nil {
				item.Error = errObj.Error
			} else {
				item.Error = "unknown refill response"
			}
		}
		items = append(items, item)
	}

	return smmprams.MultiRefillResponse{Items: items}, nil
}
