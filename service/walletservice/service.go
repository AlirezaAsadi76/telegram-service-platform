package walletservice

type Service struct {
	repo            Repository
	txRepo          TransactionRepository
	idempotencyRepo IdempotencyChecker
	config          Config
}

func New(repo Repository, txRepo TransactionRepository, idempotencyRepo IdempotencyChecker, cfg Config) *Service {
	return &Service{
		repo:            repo,
		txRepo:          txRepo,
		idempotencyRepo: idempotencyRepo,
		config:          cfg,
	}
}
