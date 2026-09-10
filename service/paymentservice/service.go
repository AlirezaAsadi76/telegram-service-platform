package paymentservice

type Service struct {
	repo                    Repository
	paymentConfirmationRepo PaymentConfirmationRepository
	zarinpalAdapter         Provider
	cryptoAdapter           Provider
}

func New(
	repo Repository,
	paymentConfirmationRepo PaymentConfirmationRepository,
	zarinpal Provider,
	crypto Provider,
) *Service {
	return &Service{
		repo:                    repo,
		paymentConfirmationRepo: paymentConfirmationRepo,
		zarinpalAdapter:         zarinpal,
		cryptoAdapter:           crypto,
	}
}
