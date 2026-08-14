package notificationdispatchjob

import (
	"context"
	"errors"
	"fmt"
	"log"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/params/notificationparams"
	"telegram-service-platform/pkg/unmarshal"

	"github.com/redis/go-redis/v9"
)

func (j *Job) Run(ctx context.Context) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	var notification notificationentity.Notification
	fromQueue := false

	// ۱. ابتدا Redis Queue
	result, brErr := j.redis.BRPop(ctx, j.config.queueKey, j.config.timeout)

	if brErr == nil && len(result) >= 2 {
		if notify, unmarshalErr := unmarshal.UnmarshalToNotification(result[1]); unmarshalErr == nil {
			fromQueue = true
			notification = notify
		}
	} else if brErr != nil && !errors.Is(brErr, redis.Nil) {
		return brErr
	}

	// ۲. اگر Queue خالی بود، از DB پولینگ کن (Retry logic)
	if !fromQueue {
		pending, err := j.notificationService.GetPending(ctx, notificationparams.GetPendingRequest{
			Limit: j.config.PendingLimit,
		})
		if err != nil {
			return err
		}
		if len(pending.Notifications) == 0 {
			return nil
		}
		notification = pending.Notifications[0]
	}

	text := j.buildMessage(notification)

	if err := j.bot.SendText(ctx, int64(notification.UserID), text); err != nil {
		log.Printf("send notification %d failed: %v", notification.ID, err)

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
	return nil
}

func (j *Job) buildMessage(n notificationentity.Notification) string {
	switch n.Type {
	case notificationentity.NotificationTypeOrderPaid:
		return fmt.Sprintf("✅ سفارش شما #%v با موفقیت پرداخت شد و در حال پردازش است.", n.Payload["order_id"])
	case notificationentity.NotificationTypeOrderCompleted:
		return fmt.Sprintf("🎉 سفارش شما #%v تکمیل شد!", n.Payload["order_id"])
	case notificationentity.NotificationTypeOrderFailed:
		return fmt.Sprintf("❌ سفارش شما #%v ناموفق بود. لطفاً با پشتیبانی تماس بگیرید.", n.Payload["order_id"])
	case notificationentity.NotificationTypePaymentExpired:
		return fmt.Sprintf("⏳ زمان پرداخت سفارش #%v به پایان رسید. لطفاً دوباره تلاش کنید.", n.Payload["order_id"])
	case notificationentity.NotificationTypeWalletRecharged:
		return fmt.Sprintf("💰 کیف پول شما به مبلغ %v ریال شارژ شد.", n.Payload["amount"])
	default:
		return fmt.Sprintf("🔔 پیام جدید: %v", n.Payload)
	}
}
