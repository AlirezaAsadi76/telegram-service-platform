package rediscatalog

import "time"

type Config struct {
	PlatformsCacheKey  string        `koanf:"platformsCacheKey"`
	categoriesCacheKey string        `koanf:"categoriesCacheKey"`
	cacheTTL           time.Duration `koanf:"cacheTTL"`
}
