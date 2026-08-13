package paymentservice

import "telegram-service-platform/entity/paymententity"

func (s *Service) getProvider(method paymententity.PaymentMethod) Provider {
	switch method {
	case paymententity.PaymentMethodZarinpal:
		return s.zarinpalAdapter
	case paymententity.PaymentMethodCrypto:
		return s.cryptoAdapter
	default:
		return s.zarinpalAdapter
	}
}
