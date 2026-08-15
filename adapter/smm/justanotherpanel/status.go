package justanotherpanel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"telegram-service-platform/entity/smmentity"
)

func (a *Adapter) Status(ctx context.Context) error {
	form := url.Values{
		"key":    {a.config.APIKey},
		"action": {string(ActionTypeStatus)},
		"order":  {"996993167"},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		a.config.BaseURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"JAP API returned %s",
			resp.Status,
		)
	}
	var statusOrder smmentity.Status
	if dErr := json.NewDecoder(resp.Body).Decode(&statusOrder); dErr != nil {
		return dErr
	}
	fmt.Println("Status:", statusOrder)
	return nil
}
