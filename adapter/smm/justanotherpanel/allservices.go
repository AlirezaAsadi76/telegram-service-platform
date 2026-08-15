package justanotherpanel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/params/smmprams"
)

func (a *Adapter) AllServices(ctx context.Context) (smmprams.GetAllServicesResponse, error) {
	form := url.Values{
		"key":    {a.config.APIKey},
		"action": {string(ActionTypeServices)},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		a.config.BaseURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return smmprams.GetAllServicesResponse{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return smmprams.GetAllServicesResponse{}, err

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return smmprams.GetAllServicesResponse{}, fmt.Errorf(
			"JAP API returned %s",
			resp.Status,
		)
	}

	Smms := make([]smmentity.SMM, 0)
	if err := json.NewDecoder(resp.Body).Decode(&Smms); err != nil {
		return smmprams.GetAllServicesResponse{}, err

	}

	return smmprams.GetAllServicesResponse{
		Services: Smms,
	}, nil
}
