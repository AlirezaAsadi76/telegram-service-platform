package smmproviderservice

import "time"

func (s *Service) RegisterProvider(name string, adapter SMMProvider) {
	s.providers[name] = adapter
	s.breakers[name] = NewCircuitBreaker(5, 3, 30*time.Second)
}
