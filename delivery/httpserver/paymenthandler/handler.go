package paymenthandler

type Handler struct {
	paymentService PaymentFlow
	paymentVal     PaymentValidator
	middleware     AuthMiddleware
}

func New(
	paymentService PaymentFlow,
	paymentVal PaymentValidator,
	middleware AuthMiddleware,
) *Handler {
	return &Handler{
		paymentService: paymentService,
		paymentVal:     paymentVal,
		middleware:     middleware,
	}
}
