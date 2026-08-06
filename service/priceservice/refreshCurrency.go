package priceservice

import "context"

func (s Service) RefreshCurrency(ctx context.Context) error {

	tonUsd, tErr := s.currency.GetTonUsdPrice(ctx)
	if tErr != nil {
		return tErr
	}
	usdIrr, uErr := s.currency.GetUsdTomanPrice(ctx)
	if uErr != nil {
		return uErr
	}
	stErr := s.repository.SetTonUsdPrice(ctx, tonUsd)
	if stErr != nil {
		return stErr
	}
	suErr := s.repository.SetUsdTomanPrice(ctx, usdIrr)
	if suErr != nil {
		return suErr

	}

	return nil
}
