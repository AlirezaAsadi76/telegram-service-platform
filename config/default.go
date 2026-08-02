package config

import "time"

var defaultValue = map[string]any{
	"auth.access_subject":                   AccessTokenSubject,
	"auth.refresh_subject":                  RefreshTokenSubject,
	"auth.access_token_duration":            AccessTokenDuration,
	"auth.refresh_token_duration":           RefreshTokenDuration,
	"application.graceful_shutdown_timeout": 5 * time.Second,
}
