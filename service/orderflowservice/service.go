package orderflowservice

type Service struct {
	repo   Repository
	config Config
}

func New(repo Repository, config Config) *Service {
	return &Service{
		repo:   repo,
		config: config,
	}
}
