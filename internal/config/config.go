package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const envAPIKey = "WEATHER_API_KEY"

// Default OpenWeatherMap API base URL.
const defaultAPIBaseURL = "https://api.openweathermap.org/data/2.5"

// Config holds application configuration.
type Config struct {
	APIKey     string `json:"api_key"`
	APIBaseURL string `json:"api_base_url,omitempty"`
}

// Load reads configuration from the environment (and optionally overrides from file).
// WEATHER_API_KEY must be set. API base URL defaults if not set.
func Load() (*Config, error) {
	key := os.Getenv(envAPIKey)
	if key == "" {
		return nil, fmt.Errorf("missing %s", envAPIKey)
	}
	if !ValidateAPIKey(key) {
		return nil, fmt.Errorf("invalid API key: must be 32 hex characters")
	}
	return &Config{
		APIKey:     key,
		APIBaseURL: defaultOrEnv("WEATHER_API_URL", defaultAPIBaseURL),
	}, nil
}

// LoadFromFile reads configuration from a JSON file. File format: {"api_key": "..."}.
// API key is validated.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}
	if c.APIBaseURL == "" {
		c.APIBaseURL = defaultAPIBaseURL
	}
	if !ValidateAPIKey(c.APIKey) {
		return nil, fmt.Errorf("invalid API key in config: must be 32 hex characters")
	}
	return &c, nil
}

// ValidateAPIKey returns true if key is non-empty and 32 characters (OWM format).
func ValidateAPIKey(key string) bool {
	if key == "" || len(key) != 32 {
		return false
	}
	for _, r := range key {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
			return false
		}
	}
	return true
}

func defaultOrEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
