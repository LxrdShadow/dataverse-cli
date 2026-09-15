package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	BaseURL string
	Token   string
	Timeout time.Duration
}

var DEFAULT_TIMEOUT = 1 * time.Second

func Load() (Config, error) {
	// Get the base URL and token from environment variables
	baseURL := os.Getenv("DVC_URL")
	token := os.Getenv("DVC_TOKEN")
	if baseURL == "" || token == "" {
		return Config{}, fmt.Errorf("DVC_URL and DVC_TOKEN must be set")
	}

	return Config{
		BaseURL: baseURL,
		Token:   token,
		Timeout: DEFAULT_TIMEOUT,
	}, nil
}
