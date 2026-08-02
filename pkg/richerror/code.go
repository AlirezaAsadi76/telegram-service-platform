package richerror

type Code string

const (
	CodeUnknown Code = "UNKNOWN"

	CodeUserNotFound Code = "USER_NOT_FOUND"

	CodeUserAlreadyExists Code = "USER_ALREADY_EXISTS"

	CodeProductNotFound Code = "PRODUCT_NOT_FOUND"

	CodeOrderNotFound Code = "ORDER_NOT_FOUND"

	CodeOrderAlreadyPaid Code = "ORDER_ALREADY_PAID"

	CodePaymentFailed Code = "PAYMENT_FAILED"

	CodePaymentAlreadyProcessed Code = "PAYMENT_ALREADY_PROCESSED"

	CodeTelegramOperationFailed Code = "TELEGRAM_OPERATION_FAILED"

	CodeInvalidInput Code = "INVALID_INPUT"
)
