package smmproviderservice

type Service struct {
	repo      Repository
	providers map[string]SMMProvider // name -> adapter
	breakers  map[string]*CircuitBreaker
}

func New(repo Repository) *Service {
	return &Service{
		repo:      repo,
		providers: make(map[string]SMMProvider),
		breakers:  make(map[string]*CircuitBreaker),
	}
}
