package checkoutservice

type Config struct {
	PrefixDirectIdempotencyKey string `koanf:"prefix_idempotency_key"`
	PrefixWalletIdempotencyKey string `koanf:"prefix_wallet_idempotency_key"`
	PrefixManualIdempotencyKey string `koanf:"prefix_Manual_idempotency_key"`
}
