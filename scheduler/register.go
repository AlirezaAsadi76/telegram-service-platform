package scheduler

import (
	"log"

	"github.com/go-co-op/gocron/v2"
)

func (s *Scheduler) Register() error {

	for _, job := range s.jobs {

		_, err := s.engine.NewJob(gocron.DurationJob(s.config.CurrencyRefreshInterval),
			gocron.NewTask(
				func() {

					if err := job.Run(s.ctx); err != nil {
						log.Printf(
							"job %s failed: %v",
							job.Name(),
							err,
						)
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
