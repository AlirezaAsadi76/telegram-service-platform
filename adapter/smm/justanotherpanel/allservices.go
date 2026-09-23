package justanotherpanel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/params/smmparams"
)

func (a *Adapter) AllServices(ctx context.Context) (smmparams.GetAllServicesResponse, error) {
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
		return smmparams.GetAllServicesResponse{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, dErr := a.client.Do(req)
	if dErr != nil {
		return smmparams.GetAllServicesResponse{}, dErr

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return smmparams.GetAllServicesResponse{}, fmt.Errorf(
			"JAP API returned %s",
			resp.Status,
		)
	}

	services := make([]smmentity.SMM, 0)
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		return smmparams.GetAllServicesResponse{}, err

	}

	for i := range services {
		services[i].IsActive = true
		services[i].ProviderName = ProviderName
	}

	return smmparams.GetAllServicesResponse{
		Services: services,
	}, nil
}
