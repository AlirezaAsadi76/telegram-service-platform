package paymentparams

import "telegram-service-platform/entity/paymententity"

type UpdateStatusRequest struct {
	PymentId uint64
	Status   paymententity.PaymentStatus
}
