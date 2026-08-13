package smmproviderservice

func (s *Service) RegisterProvider(name string, adapter SMMProvider) {
	s.providers[name] = adapter
	s.breakers[name] = NewCircuitBreaker(s.config.FailureThreshold, s.config.SuccessThreshold, s.config.circuitTimeout)
}
