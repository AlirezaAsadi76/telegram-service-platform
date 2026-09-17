package msgerror

const (
	InternalServerError          = "something went wrong"
	IdempotencyAlreadyProcessing = "idempotency already processing"
	InvalidInput                 = "invalid input"
	Unexpected                   = "unexpected"
	QueryFailed                  = "failed to query"
	QueryScanFailed              = "failed to query scan"
	ExternalServiceFailed        = "external service failed"
	InvalidPrice                 = "invalid price"
	MarshalFailed                = "marshal failed"
	UnmarshalFailed              = "unmarshal failed"
	NoAvailableAdapter           = "no available adapter"
	InsufficientBalance          = "insufficient balance"
)

const (
	CacheEmpty       = "cache empty"
	CacheNotFound    = "cache not found"
	CacheReadFailed  = "failed to read cache"
	CacheWriteFailed = "failed to write cache"
	CacheParseFailed = "failed to parse cache value"
)

const (
	OrderCreateFailed             = "order create failed"
	OrderUpdateFailed             = "order update failed"
	OrderNotFound                 = "order not found"
	OrderFulfillmentEnqueueFailed = "failed to enqueue order fulfillment"
	OrderAlreadyPaid              = "order already paid"
)

const (
	SMMProviderRejected        = "smm provider rejected the order"
	SMMProviderUnavailable     = "smm provider unavailable"
	SMMProviderTimeout         = "smm provider timeout"
	SMMProviderInvalidResponse = "smm provider returned an invalid response"
	SMMProviderHTTPError       = "smm provider http error"
	SMMProviderRequestFailed   = "failed to call smm provider"
)

const (
	ErrInvalidToken     = "invalid or malformed token"
	ErrTokenExpired     = "token has expired"
	ErrInvalidAlgorithm = "invalid signing algorithm"
	ErrMissingBearer    = "missing Bearer prefix"
)

const (
	PaymentConfirmationConflict = "payment confirmation conflict"
	PaymentProviderError        = "payment provider error"
	PaymentVerifyFailed         = "payment verify failed"
	PaymentAlreadyProcessed     = "payment already processed"
	PaymentNotFound             = "payment not found"
	PaymentFailed               = "payment failed"
)

const (
	WalletTransactionNotFound  = "wallet transaction not found"
	WalletIdempotencyKeyReused = "wallet idempotency key reused"
	WalletNotFound             = "wallet not found"
)

const (
	UserNotFound      = "user not found"
	UserAlreadyExists = "user already exists"
	UserUnauthorized  = "user unauthorized"
)

const (
	ProductNotFound = "productkeyboard not found"
)
