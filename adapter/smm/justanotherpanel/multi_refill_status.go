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

// MultiRefillStatus queries the status of multiple refill requests at once.
// JAP returns: [
//
//	{"refill": 1, "status": "Completed"},
//	{"refill": 2, "status": "Rejected"},
//	{"refill": 3, "status": {"error": "Refill not found"}}
//
// ]
func (a *Adapter) MultiRefillStatus(ctx context.Context, refillIDs []string) (smmprams.MultiRefillStatusResponse, error) {
	form := url.Values{
		"key":     {a.config.APIKey},
		"action":  {string(ActionTypeRefillStatus)},
		"refills": {strings.Join(refillIDs, ",")},
	}

	// JAP returns array where "status" can be string or {"error": "..."}
	var rawItems []struct {
		Refill int64           `json:"refill"`
		Status json.RawMessage `json:"status"`
	}
	if err := a.doRequest(ctx, form, &rawItems); err != nil {
		return smmprams.MultiRefillStatusResponse{}, fmt.Errorf("multi-refill-status request: %w", err)
	}

	items := make([]smmentity.RefillStatusItem, 0, len(rawItems))
	for _, r := range rawItems {
		item := smmentity.RefillStatusItem{RefillID: r.Refill}

		// Try parse as string first (success case)
		var statusStr string
		if err := json.Unmarshal(r.Status, &statusStr); err == nil {
			item.Status = statusStr
		} else {
			// Parse as error object: {"error": "..."}
			var errObj struct {
				Error string `json:"error"`
			}
			if err := json.Unmarshal(r.Status, &errObj); err == nil {
				item.Error = errObj.Error
			} else {
				item.Error = "unknown status response"
			}
		}
		items = append(items, item)
	}

	return smmprams.MultiRefillStatusResponse{Items: items}, nil
}
