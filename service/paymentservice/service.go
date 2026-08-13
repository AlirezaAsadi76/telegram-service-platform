// service/paymentservice/service.go
package paymentservice

type Service struct {
	repo            Repository
	zarinpalAdapter Provider
	cryptoAdapter   Provider
}

func New(repo Repository, zarinpal, crypto Provider) *Service {
	return &Service{
		repo:            repo,
		zarinpalAdapter: zarinpal,
		cryptoAdapter:   crypto,
	}
}
