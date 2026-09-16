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

	resp, err := a.client.Do(req)
	if err != nil {
		return smmparams.GetAllServicesResponse{}, err

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return smmparams.GetAllServicesResponse{}, fmt.Errorf(
			"JAP API returned %s",
			resp.Status,
		)
	}

	Smms := make([]smmentity.SMM, 0)
	if err := json.NewDecoder(resp.Body).Decode(&Smms); err != nil {
		return smmparams.GetAllServicesResponse{}, err

	}

	return smmparams.GetAllServicesResponse{
		Services: Smms,
	}, nil
}
