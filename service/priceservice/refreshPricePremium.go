package priceservice

import "context"

func (s Service) RefreshPricePremium(ctx context.Context) error {
	const Op = "priceservice.RefreshPricePremium"
	starsPricePrv, gErr := s.telegramPrv.GetPremiumPlans(ctx)
	if gErr != nil {
		return gErr
	}
	sErr := s.repository.SetPremiumPrices(ctx, starsPricePrv)
	if sErr != nil {
		return sErr
	}
	return nil

}
