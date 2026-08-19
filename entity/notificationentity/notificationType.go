package notificationentity

type NotificationType string

const (
	NotificationTypeOrderPaid       NotificationType = "ORDER_PAID"
	NotificationTypeOrderCompleted  NotificationType = "ORDER_COMPLETED"
	NotificationTypeOrderFailed     NotificationType = "ORDER_FAILED"
	NotificationTypePaymentExpired  NotificationType = "PAYMENT_EXPIRED"
	NotificationTypeWalletRecharged NotificationType = "WALLET_RECHARGED"
	NotificationTypeSystemAlert     NotificationType = "SYSTEM_ALERT"
)
