package smmcatalogsyncjob

import (
	"context"
	"time"

	"telegram-service-platform/logger"
	"telegram-service-platform/pkg/metrics"

	"go.uber.org/zap"
)

func (j *Job) Run(ctx context.Context) error {
	start := time.Now()
	jobName := j.Name()

	defer func() {
		metrics.WorkerDuration.
			WithLabelValues(jobName).
			Observe(time.Since(start).Seconds())
	}()

	j.mutex.Lock()
	defer j.mutex.Unlock()

	logger.Logger.Info(
		"SMM catalog sync job started",
		zap.String("job", jobName),
	)

	if err := j.productService.SyncSMMServices(ctx); err != nil {
		logger.Logger.Error(
			"SMM catalog sync job failed",
			zap.String("job", jobName),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)

		return err
	}

	logger.Logger.Info(
		"SMM catalog sync job completed",
		zap.String("job", jobName),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}
