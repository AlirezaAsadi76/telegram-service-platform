// scheduler/register.go
package scheduler

import (
	"context"
	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/metrics"
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

func (s *Scheduler) Register() error {
	logger.Logger.Info("scheduler register start")

	for _, job := range s.jobs {
		interval := s.resolveInterval(job.Name())
		j := job

		_, err := s.engine.NewJob(
			gocron.DurationJob(interval),
			gocron.NewTask(
				func() {
					ctx := context.Background()
					start := time.Now()
					jobName := j.Name()

					if err := j.Run(ctx); err != nil {
						metrics.WorkerRuns.WithLabelValues(jobName, "error").Inc()
						logger.Logger.Error("job execution failed",
							zap.String("job", jobName),
							zap.Error(err),
							zap.Duration("duration", time.Since(start)),
						)
					} else {
						metrics.WorkerRuns.WithLabelValues(jobName, "success").Inc()
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
		logger.Logger.Info("job registered",
			zap.String("job", j.Name()),
			zap.Duration("interval", interval),
		)
	}

	return nil
}

func (s *Scheduler) resolveInterval(name string) time.Duration {
	switch name {
	case "user-activity-sync":
		return s.config.userActivitysyncInterval
	case "smm-validation":
		return s.config.SmmValidationInterval
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
