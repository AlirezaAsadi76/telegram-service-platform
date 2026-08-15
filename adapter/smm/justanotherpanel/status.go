package justanotherpanel

import (
	"context"
	"fmt"
	"net/url"

	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/params/smmprams"
)

// Status queries the status of a single order from JAP API.
func (a *Adapter) Status(ctx context.Context, orderID string) (smmprams.GetStatusResponse, error) {
	form := url.Values{
		"key":    {a.config.APIKey},
		"action": {string(ActionTypeStatus)},
		"order":  {orderID},
	}

	var status smmentity.Status
	if err := a.doRequest(ctx, form, &status); err != nil {
		return smmprams.GetStatusResponse{}, fmt.Errorf("status request: %w", err)
	}

	return smmprams.GetStatusResponse{Status: status}, nil
}
