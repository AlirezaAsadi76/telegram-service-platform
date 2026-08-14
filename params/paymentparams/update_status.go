package paymentparams

import "telegram-service-platform/entity/paymententity"

type UpdateStatusRequest struct {
	PaymentId uint64
	Status    paymententity.PaymentStatus
}
