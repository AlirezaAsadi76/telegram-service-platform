package walletservice

type Service struct {
	repo            Repository
	txRepo          TransactionRepository
	idempotencyRepo IdempotencyChecker
}

func New(repo Repository, txRepo TransactionRepository, idempotencyRepo IdempotencyChecker) *Service {
	return &Service{
		repo:            repo,
		txRepo:          txRepo,
		idempotencyRepo: idempotencyRepo,
	}
}
