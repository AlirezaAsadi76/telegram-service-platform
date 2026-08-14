package paymentparams

import "telegram-service-platform/entity/paymententity"

type PendingResponse struct {
	Payments []paymententity.Payment
}
