package pricerefreshjob

import (
	"context"
	"log"
)

func (j *Job) Run(ctx context.Context) error {
	log.Println("job run")

	j.mutex.Lock()
	defer j.mutex.Unlock()

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
