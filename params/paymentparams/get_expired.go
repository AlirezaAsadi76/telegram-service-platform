package paymentparams

import "telegram-service-platform/entity/paymententity"

type ExpiredResponse struct {
	Payments []paymententity.Payment
}
