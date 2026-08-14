// scheduler/register.go
package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func (s *Scheduler) Register() error {
	fmt.Println("scheduler register start")

	for _, job := range s.jobs {
		interval := s.resolveInterval(job.Name())

		_, err := s.engine.NewJob(
			gocron.DurationJob(interval),
			gocron.NewTask(
				func() {
					ctx := context.Background()
					if err := job.Run(ctx); err != nil {
						log.Printf("job %s failed: %v", job.Name(), err)
					}
				},
			),
			gocron.WithStartAt(
				gocron.WithStartImmediately(),
			),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) resolveInterval(name string) time.Duration {
	switch name {
	case "price-refresh":
		return s.config.CurrencyRefreshInterval
	case "payment-verify":
		return s.config.PaymentVerifyInterval
	case "status-sync":
		return s.config.StatusSyncInterval
	case "payment-expiry":
		return s.config.PaymentExpiryInterval
	case "order-fulfiller", "notification-dispatch":
		return s.config.QueueConsumerInterval
	default:
		return 1 * time.Minute
	}
}
