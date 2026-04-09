package config

import (
	"os"
	"time"
)

type Config struct {
	AppName           string
	AppEnv            string
	AppPort           string
	HTTPClientTimeout time.Duration
	CatalogCacheTTL   time.Duration
	StreamCacheTTL    time.Duration
	LogLevel          string
}

func Load() *Config {
	return &Config{
		AppName:           getEnv("APP_NAME", "cameras-api"),
		AppEnv:            getEnv("APP_ENV", "local"),
		AppPort:           getEnv("APP_PORT", "8080"),
		HTTPClientTimeout: parseDuration(getEnv("HTTP_CLIENT_TIMEOUT", "10s")),
		CatalogCacheTTL:   parseDuration(getEnv("CATALOG_CACHE_TTL", "15m")),
		StreamCacheTTL:    parseDuration(getEnv("STREAM_CACHE_TTL", "1m")),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 10 * time.Second
	}
	return d
}
