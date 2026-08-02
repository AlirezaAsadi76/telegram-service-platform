package config

import "time"

const (
	jwtSecret                = "secret"
	AccessTokenSubject       = "at"
	RefreshTokenSubject      = "rt"
	AccessTokenDuration      = time.Hour * 24
	RefreshTokenDuration     = time.Hour * 24 * 7
	AuthMiddlewareContextKey = "claims"
)
