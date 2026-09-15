package paymentservice

type Service struct {
	repo                    Repository
	paymentConfirmationRepo PaymentConfirmationRepository
	fulfillmentEnqueuer     FulfillmentEnqueuer
	zarinpalAdapter         Provider
	cryptoAdapter           Provider
}

func New(
	repo Repository,
	paymentConfirmationRepo PaymentConfirmationRepository,
	zarinpal Provider,
	crypto Provider,
	fulfillmentEnqueuer ...FulfillmentEnqueuer,
) *Service {
	var enqueuer FulfillmentEnqueuer
	if len(fulfillmentEnqueuer) > 0 {
		enqueuer = fulfillmentEnqueuer[0]
	}
	return &Service{
		repo:                    repo,
		paymentConfirmationRepo: paymentConfirmationRepo,
		fulfillmentEnqueuer:     enqueuer,
		zarinpalAdapter:         zarinpal,
		cryptoAdapter:           crypto,
	}
}
