package useractivitysyncjob

import (
	"context"
	"time"

	"go.uber.org/zap"

	"telegram-service-platform/logger"
	"telegram-service-platform/params"
	"telegram-service-platform/pkg/metrics"
)

func (j *Job) Run(ctx context.Context) error {
	start := time.Now()
	jobName := j.Name()

	defer func() {
		metrics.WorkerDuration.WithLabelValues(jobName).Observe(time.Since(start).Seconds())
	}()

	logger.Logger.Info("user activity sync job started", zap.String("job", jobName))

	resp, err := j.userService.SyncActiveUsersLastSeen(ctx, params.SyncLastSeenRequest{})
	if err != nil {
		metrics.WorkerRuns.WithLabelValues(jobName, "error").Inc()
		logger.Logger.Error("user activity sync failed", zap.String("job", jobName), zap.Error(err))
		return err
	}

	metrics.WorkerRuns.WithLabelValues(jobName, "success").Inc()
	logger.Logger.Info("user activity sync completed",
		zap.String("job", jobName),
		zap.Int64("rows_updated", resp.Synced),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}
