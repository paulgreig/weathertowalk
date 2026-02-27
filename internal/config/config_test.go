package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_APIKeyFromEnvironment(t *testing.T) {
	key := "abc123def456789012345678901234ab"
	os.Setenv(envAPIKey, key)
	defer os.Unsetenv(envAPIKey)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() err = %v, want nil", err)
	}
	if cfg.APIKey != key {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, key)
	}
}

func TestLoad_MissingAPIKeyReturnsError(t *testing.T) {
	os.Unsetenv(envAPIKey)

	_, err := Load()
	if err == nil {
		t.Fatal("Load() err = nil, want error when API key missing")
	}
}

func TestLoad_InvalidAPIKeyReturnsError(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"empty", ""},
		{"too_short", "short"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(envAPIKey, tt.key)
			defer os.Unsetenv(envAPIKey)

			_, err := Load()
			if err == nil {
				t.Errorf("Load() err = nil, want error for invalid key %q", tt.key)
			}
		})
	}
}

func TestLoad_DefaultAPIURL(t *testing.T) {
	os.Setenv(envAPIKey, "a1b2c3d4e5f6789012345678901234ab")
	defer os.Unsetenv(envAPIKey)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() err = %v, want nil", err)
	}
	if cfg.APIBaseURL == "" {
		t.Error("APIBaseURL is empty, want default")
	}
}

func TestLoad_APIKeyFromConfigFileWhenEnvUnset(t *testing.T) {
	os.Unsetenv(envAPIKey)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{"api_key": "f0e1d2c3b4a5987654321098765432ab"}`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile() err = %v, want nil", err)
	}
	if cfg.APIKey != "f0e1d2c3b4a5987654321098765432ab" {
		t.Errorf("APIKey = %q, want from file", cfg.APIKey)
	}
}

// TestLoadFromFile_DefaultAPIBaseURLWhenOmitted ensures that when the config file
// omits api_base_url, LoadFromFile sets it to the default. (Mutation testing:
// survived mutant was "if c.APIBaseURL == ”" changed to "!= ”", skipping default.)
func TestLoadFromFile_DefaultAPIBaseURLWhenOmitted(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{"api_key": "a1b2c3d4e5f6789012345678901234ab"}`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile() err = %v, want nil", err)
	}
	if cfg.APIBaseURL == "" {
		t.Error("APIBaseURL is empty; want default when omitted from file")
	}
	if cfg.APIBaseURL != "https://api.openweathermap.org/data/2.5" {
		t.Errorf("APIBaseURL = %q, want default URL", cfg.APIBaseURL)
	}
}

func TestValidateAPIKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want bool
	}{
		{"valid_32", "a1b2c3d4e5f6789012345678901234ab", true},
		{"empty", "", false},
		{"too_short", "abc", false},
		{"31_chars", "a1b2c3d4e5f6789012345678901234a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateAPIKey(tt.key)
			if got != tt.want {
				t.Errorf("ValidateAPIKey(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}
