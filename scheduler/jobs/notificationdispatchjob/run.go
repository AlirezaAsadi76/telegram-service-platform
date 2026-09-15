package notificationdispatchjob

import (
	"context"
	"errors"
	"fmt"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/pkg/metrics"
	"telegram-service-platform/pkg/unmarshal"
	"time"

	"github.com/go-telegram/bot"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (j *Job) Run(ctx context.Context) error {
	start := time.Now()
	jobName := j.Name()

	defer func() {
		metrics.WorkerDuration.WithLabelValues(jobName).Observe(time.Since(start).Seconds())
	}()

	logger.Logger.Info("worker started", zap.String("job", jobName))

	j.mutex.Lock()
	defer j.mutex.Unlock()

	var notification notificationentity.Notification
	fromQueue := false

	// ۱. ابتدا Redis Queue
	result, brErr := j.redis.BRPop(ctx, j.config.Timeout, j.config.QueueKey)

	if brErr == nil && len(result) >= 2 {
		if notify, unmarshalErr := unmarshal.UnmarshalToNotification(result[1]); unmarshalErr == nil {
			fromQueue = true
			notification = notify
		}
	} else if brErr != nil && !errors.Is(brErr, redis.Nil) {

		metrics.WorkerRuns.WithLabelValues(jobName, "error").Inc()
		logger.Logger.Error("worker redis error", zap.String("job", jobName), zap.Error(brErr))
		return brErr
	}

	// ۲. اگر Queue خالی بود، از DB پولینگ کن (Retry logic)
	if !fromQueue {
		pending, err := j.notificationService.GetPending(ctx, notificationparams.GetPendingRequest{
			Limit: j.config.PendingLimit,
		})
		if err != nil {
			metrics.WorkerRuns.WithLabelValues(jobName, "error").Inc()
			logger.Logger.Error("worker getPending error", zap.String("job", jobName), zap.Error(err))
			return err
		}
		if len(pending.Notifications) == 0 {
			return nil
		}
		notification = pending.Notifications[0]
	}

	text := j.buildMessage(notification)

	if err := j.bot.Send(ctx, &bot.SendMessageParams{
		ChatID: int64(notification.UserID),
		Text:   text,
	}); err != nil {

		logger.Logger.Warn("send notification failed",
			zap.Uint64("notification_id", notification.ID),
			zap.Error(err),
		)

		if notification.RetryCount >= 2 {
			_ = j.notificationService.UpdateStatus(ctx, notificationparams.UpdateStatusRequest{
				Id:     notification.ID,
				Status: notificationentity.NotificationStatusFailed,
			})
		} else {
			_ = j.notificationService.IncrementRetry(ctx, notificationparams.IncrementRetryRequest{
				Id:         notification.ID,
				RetryCount: notification.RetryCount + 1,
			})
		}
		return nil
	}

	_ = j.notificationService.UpdateStatus(ctx, notificationparams.UpdateStatusRequest{
		Id:     notification.ID,
		Status: notificationentity.NotificationStatusSent,
	})

	metrics.NotificationsSent.WithLabelValues("success").Inc()
	logger.Logger.Info("notification sent",
		zap.String("job", jobName),
		zap.Uint64("notification_id", notification.ID),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}
