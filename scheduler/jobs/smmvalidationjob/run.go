package smmvalidationjob

import (
	"context"
	"fmt"
	"strings"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/logger"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/pkg/metrics"
	"time"

	"go.uber.org/zap"
)

func (j *Job) Run(ctx context.Context) error {
	start := time.Now()
	jobName := j.Name()

	defer func() {
		metrics.WorkerDuration.WithLabelValues(jobName).Observe(time.Since(start).Seconds())
	}()

	logger.Logger.Info("SMM validation job started", zap.String("job", jobName))

	missingServices, err := j.productService.GetMissingSMMServices(ctx)
	if err != nil {
		metrics.WorkerRuns.WithLabelValues(jobName, "error").Inc()
		logger.Logger.Error("SMM validation failed", zap.String("job", jobName), zap.Error(err))
		return err
	}

	if len(missingServices) > 0 {
		metrics.WorkerRuns.WithLabelValues(jobName, "warning").Inc()

		var missingNames []string
		for _, svc := range missingServices {
			missingNames = append(missingNames, fmt.Sprintf("- %s (ID: %d, Category: %s)", svc.Name, svc.Service, svc.Category))
		}

		messageBody := fmt.Sprintf("⚠️ هشدار اعتبارسنجی SMM:\nتعداد %d سرویس در دیتابیس وجود دارد که دیگر در پرووایدر فعال نیستند.\n\n%s\n\nلطفاً از طریق پنل ادمین اقدام به غیرفعال‌سازی یا حذف آن‌ها کنید.", len(missingServices), strings.Join(missingNames, "\n"))

		// ثبت نوتیفیکیشن برای ادمین
		//TODO- use admins id in database
		notifyErr := j.notificationSvc.Create(ctx, notificationparams.CreateRequest{
			UserID: uint64(j.config.AdminUserID),
			Type:   notificationentity.NotificationTypeSystemAlert,
			Payload: map[string]any{
				"missing_count":   len(missingServices),
				"action_required": "review_smm_services",
				"message":         messageBody,
			},
		})

		if notifyErr != nil {
			logger.Logger.Error("failed to send admin notification for SMM validation", zap.Error(notifyErr))
		} else {
			logger.Logger.Warn("SMM validation found missing services, admin notified",
				zap.Int("missing_count", len(missingServices)),
				zap.String("job", jobName),
			)
		}
	} else {
		metrics.WorkerRuns.WithLabelValues(jobName, "success").Inc()
		logger.Logger.Info("SMM validation completed successfully, all services are valid", zap.String("job", jobName))
	}

	return nil
}
