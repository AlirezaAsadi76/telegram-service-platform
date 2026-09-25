package checkoutservice

import "time"

type Config struct {
	PrefixDirectIdempotencyKey string `koanf:"prefix_idempotency_key"`
	PrefixWalletIdempotencyKey string `koanf:"prefix_wallet_idempotency_key"`
	PrefixManualIdempotencyKey string `koanf:"prefix_Manual_idempotency_key"`

	WalletPriceLockTTL   time.Duration `koanf:"wallet_price_lock_ttl"`
	ZarinpalPriceLockTTL time.Duration `koanf:"zarinpal_price_lock_ttl"`
	CryptoPriceLockTTL   time.Duration `koanf:"crypto_price_lock_ttl"`
}
