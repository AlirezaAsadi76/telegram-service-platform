package priceservice

import "context"

func (s Service) RefreshPriceStars(ctx context.Context) error {
	const Op = "priceservice.refreshPriceStars"
	starsPricePrv, gErr := s.telegramPrv.GetStarPrice(ctx)
	if gErr != nil {
		return gErr
	}
	sErr := s.repository.SetStarPrice(ctx, starsPricePrv)
	if sErr != nil {
		return sErr
	}
	return nil

}
