package authservice

import "time"

type Config struct {
	SignKey                  string        `koanf:"sign_key"`
	AccessTokenDuration      time.Duration `koanf:"access_token_duration"`
	RefreshTokenDuration     time.Duration `koanf:"refresh_token_duration"`
	AccessSubject            string        `koanf:"access_subject"`
	RefreshSubject           string        `koanf:"refresh_subject"`
	AuthMiddlewareContextKey string        `koanf:"auth_middleware_context_key"`
}
