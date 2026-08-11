package paymentservice

import "telegram-service-platform/entity/paymententity"

func (s *Service) getProvider(method paymententity.PaymentMethod) Provider {
	switch method {
	case paymententity.PaymentMethodZarinpal:
		return s.zarinpalProvider
	case paymententity.PaymentMethodCrypto:
		return s.cryptoProvider
	default:
		return s.zarinpalProvider
	}
}
