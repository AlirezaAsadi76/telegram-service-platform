package pricerefreshjob

import "context"

func (j Job) Run(ctx context.Context) error {

	err := j.priceService.RefreshCurrency(ctx)
	if err != nil {
		return err
	}

	err = j.priceService.RefreshPriceStars(ctx)
	if err != nil {
		return err
	}

	err = j.priceService.RefreshPricePremium(ctx)
	if err != nil {
		return err
	}

	return nil
}
