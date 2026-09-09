package richerror

type Code string

const (
	CodeUnknown Code = "UNKNOWN"

	CodeUserNotFound      Code = "USER_NOT_FOUND"
	CodeUserAlreadyExists Code = "USER_ALREADY_EXISTS"
	CodeProductNotFound   Code = "PRODUCT_NOT_FOUND"

	CodeOrderNotFound          Code = "ORDER_NOT_FOUND"
	CodeOrderAlreadyPaid       Code = "ORDER_ALREADY_PAID"
	CodeOrderInvalidState      Code = "ORDER_INVALID_STATE"
	CodeOrderInvalidTransition Code = "ORDER_INVALID_TRANSITION"

	CodePaymentNotFound            Code = "PAYMENT_NOT_FOUND"
	CodePaymentFailed              Code = "PAYMENT_FAILED"
	CodePaymentAlreadyProcessed    Code = "PAYMENT_ALREADY_PROCESSED"
	CodePaymentInvalidState        Code = "PAYMENT_INVALID_STATE"
	CodePaymentAlreadyConfirmed    Code = "PAYMENT_ALREADY_CONFIRMED"
	CodePaymentProviderRejected    Code = "PAYMENT_PROVIDER_REJECTED"
	CodePaymentProviderTimeout     Code = "PAYMENT_PROVIDER_TIMEOUT"
	CodePaymentProviderUnavailable Code = "PAYMENT_PROVIDER_UNAVAILABLE"
	CodePaymentVerificationFailed  Code = "PAYMENT_VERIFICATION_FAILED"
	CodePaymentCreationFailed      Code = "PAYMENT_CREATION_FAILED"

	CodePaymentInvalidAmount           Code = "PAYMENT_INVALID_AMOUNT"
	CodePaymentInvalidIdempotencyKey   Code = "PAYMENT_INVALID_IDEMPOTENCY_KEY"
	CodePaymentIdempotencyKeyReused    Code = "PAYMENT_IDEMPOTENCY_KEY_REUSED"
	CodePaymentIntentAlreadyExists     Code = "PAYMENT_INTENT_ALREADY_EXISTS"
	CodePaymentIntentCreationFailed    Code = "PAYMENT_INTENT_CREATION_FAILED"
	CodePaymentInitiationUpdateFailed  Code = "PAYMENT_INITIATION_UPDATE_FAILED"
	CodePaymentProviderInvalidResponse Code = "PAYMENT_PROVIDER_INVALID_RESPONSE"
	CodePaymentLoadFailed              Code = "PAYMENT_LOAD_FAILED"

	CodeWalletNotFound            Code = "WALLET_NOT_FOUND"
	CodeWalletInsufficientBalance Code = "WALLET_INSUFFICIENT_BALANCE"
	CodeWalletInvalidAmount       Code = "WALLET_INVALID_AMOUNT"
	CodeWalletConcurrentUpdate    Code = "WALLET_CONCURRENT_UPDATE"

	CodeInvalidInput Code = "INVALID_INPUT"

	CodeTelegramOperationFailed Code = "TELEGRAM_OPERATION_FAILED"
)
