package smmproviderservice

type Service struct {
	repo      ProviderRepository
	providers map[string]SMMProvider // name -> adapter
	breakers  map[string]*CircuitBreaker
	config    Config
}

func New(repo ProviderRepository, cfg Config) *Service {
	return &Service{
		repo:      repo,
		providers: make(map[string]SMMProvider),
		breakers:  make(map[string]*CircuitBreaker),
		config:    cfg,
	}
}
