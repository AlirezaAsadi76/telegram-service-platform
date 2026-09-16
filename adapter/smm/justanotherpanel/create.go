package justanotherpanel

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"telegram-service-platform/params/smmprams"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func (a *Adapter) Create(ctx context.Context, req smmprams.CreateOrderAdapterRequest) (smmprams.CreateOrderAdapterResponse, error) {
	const Op = "justanotherpanel.Create"

	data := url.Values{}
	data.Set("key", a.config.APIKey)
	data.Set("action", string(ActionTypeAdd))
	data.Set("service", req.ServiceID)
	data.Set("link", req.Link)
	data.Set("quantity", strconv.FormatInt(req.Quantity, 10))

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		a.config.BaseURL,
		bytes.NewBufferString(data.Encode()),
	)
	if err != nil {
		return smmprams.CreateOrderAdapterResponse{
				Outcome: smmprams.CreateOrderOutcomeUnknown,
			}, richerror.New(Op, err).
				WithKind(richerror.KindDependencyFailure).
				WithCode(richerror.CodeSMMProviderRequestFailed).
				WithMessage(msgerror.SMMProviderRequestFailed)
	}

	httpReq.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)
	httpReq.Header.Set(
		"Accept",
		"application/json",
	)

	resp, dErr := a.client.Do(httpReq)
	if dErr != nil {
		code := richerror.CodeSMMProviderUnavailable
		message := msgerror.SMMProviderUnavailable

		if errors.Is(dErr, context.DeadlineExceeded) {
			code = richerror.CodeSMMProviderTimeout
			message = msgerror.SMMProviderTimeout
		}

		return smmprams.CreateOrderAdapterResponse{
				Outcome: smmprams.CreateOrderOutcomeUnknown,
			}, richerror.New(Op, dErr).
				WithKind(richerror.KindExternalAPI).
				WithCode(code).
				WithMessage(message)
	}
	defer resp.Body.Close()

	body, rErr := io.ReadAll(resp.Body)
	if rErr != nil {
		return smmprams.CreateOrderAdapterResponse{
				Outcome: smmprams.CreateOrderOutcomeUnknown,
			}, richerror.New(Op, rErr).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodeSMMProviderInvalidResponse).
				WithMessage(msgerror.SMMProviderInvalidResponse)
	}

	var result struct {
		Order int64  `json:"order"`
		Error string `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return smmprams.CreateOrderAdapterResponse{
				Outcome: smmprams.CreateOrderOutcomeUnknown,
			}, richerror.New(Op, fmt.Errorf("invalid provider response: %w", err)).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodeSMMProviderInvalidResponse).
				WithMessage(msgerror.SMMProviderInvalidResponse)
	}

	if result.Error != "" {
		return smmprams.CreateOrderAdapterResponse{
				Outcome: smmprams.CreateOrderOutcomeRejected,
			}, richerror.New(
				Op,
				fmt.Errorf("provider rejected order: %s", result.Error),
			).
				WithKind(richerror.KindConflict).
				WithCode(richerror.CodeSMMProviderRejected).
				WithMessage(msgerror.SMMProviderRejected)
	}

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		return smmprams.CreateOrderAdapterResponse{
				Outcome: smmprams.CreateOrderOutcomeUnknown,
			}, richerror.New(Op, fmt.Errorf("provider returned HTTP status %d", resp.StatusCode)).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodeSMMProviderHTTPError).
				WithMessage(msgerror.SMMProviderHTTPError)
	}

	if result.Order <= 0 {
		return smmprams.CreateOrderAdapterResponse{
				Outcome: smmprams.CreateOrderOutcomeUnknown,
			}, richerror.New(
				Op,
				fmt.Errorf("invalid provider order id: %d", result.Order),
			).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodeSMMProviderInvalidResponse).
				WithMessage(msgerror.SMMProviderInvalidResponse)
	}

	return smmprams.CreateOrderAdapterResponse{
		Outcome:         smmprams.CreateOrderOutcomeCreated,
		ExternalOrderID: strconv.FormatInt(result.Order, 10),
	}, nil
}
