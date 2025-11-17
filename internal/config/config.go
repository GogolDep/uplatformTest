package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Host string
	Port string

	BaseCurrency string

	OXRAppID   string
	OXRBaseURL string
	OXRTimeout time.Duration

	GiphyAPIKey  string
	GiphyBaseURL string
	GiphyTimeout time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		Host:         getenv("HTTP_HOST", "0.0.0.0"),
		Port:         getenv("HTTP_PORT", "8080"),
		BaseCurrency: getenv("BASE_CURRENCY", "RUB"),

		OXRAppID:   os.Getenv("OXR_APP_ID"),
		OXRBaseURL: getenv("OXR_BASE_URL", "https://openexchangerates.org/api"),
		OXRTimeout: getenvDuration("OXR_TIMEOUT", 5*time.Second),

		GiphyAPIKey:  os.Getenv("GIPHY_API_KEY"),
		GiphyBaseURL: getenv("GIPHY_BASE_URL", "https://api.giphy.com/v1"),
		GiphyTimeout: getenvDuration("GIPHY_TIMEOUT", 5*time.Second),
	}

	if cfg.OXRAppID == "" {
		return nil, errors.New("OXR_APP_ID is required")
	}
	if cfg.GiphyAPIKey == "" {
		return nil, errors.New("GIPHY_API_KEY is required")
	}
	if cfg.BaseCurrency == "" {
		return nil, errors.New("BASE_CURRENCY is required")
	}
	return cfg, nil
}

func (c *Config) HTTPAddr() string { return fmt.Sprintf("%s:%s", c.Host, c.Port) }

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
