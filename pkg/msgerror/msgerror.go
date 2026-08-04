package msgerror

const (
	InternalServerError = "something went wrong"

	UserNotFound = "user not found"

	UserAlreadyExists = "user already exists"

	InvalidInput = "invalid input"

	Unauthorized = "unauthorized"

	Forbidden = "access forbidden"

	ProductNotFound = "product not found"

	OrderNotFound = "order not found"

	OrderAlreadyPaid = "order already paid"

	PaymentFailed = "payment failed"

	PaymentAlreadyProcessed = "payment already processed"

	TelegramOperationFailed = "telegram operation failed"
	Unexpected              = "unexpected"
)
