package paymenthandler

type Handler struct {
	paymentService PaymentConfirmer
	paymentVal     PaymentValidator
}

func New(
	paymentService PaymentConfirmer,
	paymentVal PaymentValidator,
) *Handler {
	return &Handler{
		paymentService: paymentService,
		paymentVal:     paymentVal,
	}
}
