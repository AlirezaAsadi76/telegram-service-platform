package notificationservice

type Service struct {
	repo   Repository
	redis  RedisRepository
	config Config
}

func New(repo Repository, redis RedisRepository, cfg Config) *Service {
	return &Service{
		repo:   repo,
		redis:  redis,
		config: cfg,
	}
}
