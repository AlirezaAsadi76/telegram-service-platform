package justanotherpanel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"telegram-service-platform/entity/orderentity"
	"telegram-service-platform/params/smmprams"
)

func (a *Adapter) Create(ctx context.Context, req smmprams.CreateOrderAdapterRequest) (smmprams.CreateOrderAdapterResponse, error) {
	data := url.Values{}
	data.Set("key", a.config.APIKey)
	data.Set("action", string(ActionTypeAdd))
	data.Set("service", req.ServiceID)
	data.Set("link", req.Link)
	data.Set("quantity", strconv.FormatInt(req.Quantity, 10))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, a.config.BaseURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return smmprams.CreateOrderAdapterResponse{}, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return smmprams.CreateOrderAdapterResponse{}, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return smmprams.CreateOrderAdapterResponse{}, fmt.Errorf("read body: %w", err)
	}

	// JAP returns {"order": 12345} on success, {"error": "..."} on failure
	var result struct {
		Order int64  `json:"order"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return smmprams.CreateOrderAdapterResponse{}, fmt.Errorf("unmarshal: %w, body: %s", err, string(body))
	}

	if result.Error != "" {
		return smmprams.CreateOrderAdapterResponse{}, fmt.Errorf("jap error: %s", result.Error)
	}

	return smmprams.CreateOrderAdapterResponse{
		ExternalOrderID: strconv.FormatInt(result.Order, 10),
		Status:          orderentity.OrderStatusProcessing,
	}, nil
}
