package zarinpal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"telegram-service-platform/entity"
	"telegram-service-platform/entity/paymententity"
	"telegram-service-platform/params/paymentproviderparams"
	"telegram-service-platform/pkg/richerror"
)

const verifyPath = "/pg/v4/payment/verify.json"

func (a *Adapter) Verify(ctx context.Context, req paymentproviderparams.VerifyRequest) (paymentproviderparams.VerifyResponse, error) {
	const Op = "zarinpal.verify"

	if req.ExternalID == "" {
		return paymentproviderparams.VerifyResponse{},
			richerror.New(Op, errors.New("external id is empty")).
				WithKind(richerror.KindInvalidInput).
				WithCode(richerror.CodeInvalidInput)
	}

	if req.Currency != entity.CurrencyTOMAN {
		return paymentproviderparams.VerifyResponse{},
			richerror.New(Op, fmt.Errorf("unsupported currency: %s", req.Currency)).
				WithKind(richerror.KindInvalidInput).
				WithCode(richerror.CodePaymentUnsupportedCurrency)
	}

	amount, err := req.Amount.ToInt64()
	if err != nil {
		return paymentproviderparams.VerifyResponse{}, richerror.New(Op, err).
			WithKind(richerror.KindInvalidInput).
			WithCode(richerror.CodePaymentInvalidAmount)
	}

	payload := VerifyRequest{
		MerchantID: a.config.MerchantID,
		Authority:  req.ExternalID,
		Amount:     amount,
	}

	body, mErr := json.Marshal(payload)
	if mErr != nil {
		return paymentproviderparams.VerifyResponse{},
			richerror.New(Op, mErr).
				WithKind(richerror.KindInternal)
	}

	url := strings.TrimRight(a.config.BaseURL, "/") + verifyPath

	httpReq, nrErr := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if nrErr != nil {
		return paymentproviderparams.VerifyResponse{}, richerror.New(Op, nrErr).
			WithKind(richerror.KindInternal)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, dErr := a.client.Do(httpReq)
	if dErr != nil {
		if errors.Is(dErr, context.DeadlineExceeded) ||
			errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return paymentproviderparams.VerifyResponse{}, richerror.New(Op, dErr).
				WithKind(richerror.KindExternalAPI).
				WithCode(richerror.CodePaymentProviderTimeout)
		}

		return paymentproviderparams.VerifyResponse{}, richerror.New(Op, dErr).
			WithKind(richerror.KindExternalAPI).
			WithCode(richerror.CodePaymentProviderUnavailable)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusInternalServerError {
		return paymentproviderparams.VerifyResponse{}, richerror.New(Op,
			fmt.Errorf(
				"zarinpal returned HTTP %d",
				resp.StatusCode,
			),
		).
			WithKind(richerror.KindExternalAPI).
			WithCode(richerror.CodePaymentProviderUnavailable)
	}

	var result VerifyResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return paymentproviderparams.VerifyResponse{}, richerror.New(Op, err).
			WithKind(richerror.KindExternalAPI).
			WithCode(richerror.CodePaymentProviderInvalidResponse)
	}

	if result.Errors != nil && result.Errors.Code != 0 {
		return paymentproviderparams.VerifyResponse{}, richerror.New(Op, fmt.Errorf(
			"zarinpal error %d: %s",
			result.Errors.Code,
			result.Errors.Message,
		),
		).
			WithKind(richerror.KindExternalAPI).
			WithCode(richerror.CodePaymentProviderRejected)
	}

	switch result.Data.Code {
	case 100:
		return paymentproviderparams.VerifyResponse{
			Status:      paymententity.PaymentStatusSuccess,
			ReferenceID: fmt.Sprintf("%d", result.Data.RefID),
		}, nil

	case 101:
		return paymentproviderparams.VerifyResponse{
			Status:      paymententity.PaymentStatusSuccess,
			ReferenceID: fmt.Sprintf("%d", result.Data.RefID),
		}, nil

	default:
		return paymentproviderparams.VerifyResponse{}, richerror.New(
			Op,
			fmt.Errorf(
				"zarinpal verification code: %d",
				result.Data.Code,
			),
		).
			WithKind(richerror.KindExternalAPI).
			WithCode(richerror.CodePaymentProviderRejected)
	}
}
