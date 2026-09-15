package notificationdispatchjob

import (
	"fmt"
	"telegram-service-platform/entity/notificationentity"
)

func (j *Job) buildMessage(
	n notificationentity.Notification,
) string {
	switch n.Type {
	case notificationentity.NotificationTypeOrderPaid:
		return fmt.Sprintf(
			"✅ سفارش شما #%v با موفقیت پرداخت شد و در حال پردازش است.",
			n.Payload["order_id"],
		)

	case notificationentity.NotificationTypeOrderCompleted:
		return fmt.Sprintf(
			"🎉 سفارش شما #%v تکمیل شد!",
			n.Payload["order_id"],
		)

	case notificationentity.NotificationTypeOrderFailed:
		return fmt.Sprintf(
			"❌ سفارش شما #%v ناموفق بود. لطفاً با پشتیبانی تماس بگیرید.",
			n.Payload["order_id"],
		)

	case notificationentity.NotificationTypePaymentExpired:
		return fmt.Sprintf(
			"⏳ زمان پرداخت سفارش #%v به پایان رسید. لطفاً دوباره تلاش کنید.",
			n.Payload["order_id"],
		)

	case notificationentity.NotificationTypeWalletRecharged:
		return fmt.Sprintf(
			"💰 کیف پول شما به مبلغ %v ریال شارژ شد.",
			n.Payload["amount"],
		)

	case notificationentity.NotificationTypeAdminWalletRecharge:
		return fmt.Sprintf(
			"💰 شارژ کیف پول کاربر #%v به مبلغ %v با موفقیت انجام شد.",
			n.Payload["user_id"],
			n.Payload["amount"],
		)

	default:
		return "🔔 اعلان جدیدی برای شما ثبت شده است."
	}
}
