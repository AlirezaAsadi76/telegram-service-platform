// service/paymentservice/service.go
package paymentservice

type Service struct {
	repo             Repository
	zarinpalProvider Provider
	cryptoProvider   Provider
	idempotencyRepo  IdempotencyChecker
}

func New(repo Repository, zarinpal, crypto Provider, idempotencyRepo IdempotencyChecker) *Service {
	return &Service{
		repo:             repo,
		zarinpalProvider: zarinpal,
		cryptoProvider:   crypto,
		idempotencyRepo:  idempotencyRepo,
	}
}
