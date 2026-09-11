package zarinpal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"telegram-service-platform/pkg/helpers"

	"telegram-service-platform/entity"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/pkg/richerror"
)

const createPath = "/pg/v4/payment/request.json"

func (a *Adapter) Create(ctx context.Context, req paymentproviderparams.CreateRequest) (paymentproviderparams.CreateResponse, error) {
	const Op = "zarinpal.create"

	if req.Currency != entity.CurrencyTOMAN {
		return paymentproviderparams.CreateResponse{},
			richerror.New(Op, fmt.Errorf("unsupported currency: %s", req.Currency)).
				WithKind(richerror.KindInvalidInput).
				WithCode(richerror.CodePaymentUnsupportedCurrency)
	}

	amount, err := req.Amount.ToInt64()
	if err != nil {
		return paymentproviderparams.CreateResponse{},
			richerror.New(Op, err).
				WithKind(richerror.KindInvalidInput).
				WithCode(richerror.CodePaymentInvalidAmount)
	}

	if req.CallbackURL == "" {
		return paymentproviderparams.CreateResponse{},
			richerror.New(Op, errors.New("callback url is empty")).
				WithKind(richerror.KindInvalidInput).
				WithCode(richerror.CodeInvalidInput)
	}

	payload := CreateRequest{
		MerchantID:  a.config.MerchantID,
		Amount:      amount,
		CallbackURL: req.CallbackURL,
		Description: req.Description,
	}

	body, mErr := json.Marshal(payload)
	if mErr != nil {
		return paymentproviderparams.CreateResponse{},
			richerror.New(Op, mErr).
				WithKind(richerror.KindInternal)
	}

	url := strings.TrimRight(a.config.BaseURL, "/") + createPath

	httpReq, nrErr := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if nrErr != nil {
		return paymentproviderparams.CreateResponse{},
			richerror.New(Op, nrErr).
				WithKind(richerror.KindInternal)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, dErr := a.client.Do(httpReq)
	if dErr != nil {
		if errors.Is(dErr, context.DeadlineExceeded) ||
			errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return paymentproviderparams.CreateResponse{},
				richerror.New(Op, dErr).
					WithKind(richerror.KindExternalAPI).
					WithCode(richerror.CodePaymentProviderTimeout)
		}

		return paymentproviderparams.CreateResponse{},
			richerror.New(Op, dErr).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodePaymentProviderUnavailable)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusInternalServerError {
		return paymentproviderparams.CreateResponse{},
			richerror.New(Op, fmt.Errorf("zarinpal returned HTTP %d", resp.StatusCode)).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodePaymentProviderUnavailable)
	}

	var result CreateResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return paymentproviderparams.CreateResponse{},
			richerror.New(Op, err).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodePaymentProviderInvalidResponse)
	}

	if result.Errors != nil && result.Errors.Code != 0 {
		return paymentproviderparams.CreateResponse{},
			richerror.New(Op,
				fmt.Errorf(
					"zarinpal error %d: %s",
					result.Errors.Code,
					result.Errors.Message,
				),
			).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodePaymentProviderRejected)
	}

	if result.Data.Code != 100 {
		return paymentproviderparams.CreateResponse{},
			richerror.New(Op,
				fmt.Errorf(
					"zarinpal create code: %d",
					result.Data.Code,
				),
			).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodePaymentProviderRejected)
	}

	if result.Data.Authority == "" {
		return paymentproviderparams.CreateResponse{},
			richerror.New(Op, errors.New("zarinpal returned empty authority")).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodePaymentProviderInvalidResponse)
	}

	return paymentproviderparams.CreateResponse{
		ExternalID: result.Data.Authority,
		PaymentURL: helpers.BuildPaymentURL(
			a.config.StartPayURL,
			result.Data.Authority,
		),
	}, nil
}
