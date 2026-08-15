package justanotherpanel

import (
	"context"
	"fmt"
	"net/url"

	"telegram-service-platform/params/smmprams"
)

// GetBalance returns the current account balance and currency from JAP.
func (a *Adapter) GetBalance(ctx context.Context) (smmprams.GetBalanceResponse, error) {
	form := url.Values{
		"key":    {a.config.APIKey},
		"action": {string(ActionTypeBalance)},
	}

	var result struct {
		Balance  string `json:"balance"`
		Currency string `json:"currency"`
	}
	if err := a.doRequest(ctx, form, &result); err != nil {
		return smmprams.GetBalanceResponse{}, fmt.Errorf("balance request: %w", err)
	}

	return smmprams.GetBalanceResponse{
		Balance:  result.Balance,
		Currency: result.Currency,
	}, nil
}
