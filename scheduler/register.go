package scheduler

import (
	"context"
	"log"

	"github.com/go-co-op/gocron/v2"
)

func (s *Scheduler) Register() error {

	for _, job := range s.jobs {

		_, err := s.engine.NewJob(gocron.DurationJob(s.config.CurrencyRefreshInterval),
			gocron.NewTask(
				func() {
					ctx := context.Background()

					if err := job.Run(ctx); err != nil {
						log.Println("Register job err: ", err)
						return
					}
				},
			),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
