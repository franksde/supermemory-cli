package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the CLI configuration.
type Config struct {
	APIKey              string `json:"api_key,omitempty"`
	BaseURL             string `json:"base_url,omitempty"`
	DefaultContainerTag string `json:"default_container_tag,omitempty"`
}

const DefaultBaseURL = "https://api.supermemory.ai"

// ConfigPath returns the path to the config file (~/.config/sm/config.json).
func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "sm", "config.json")
}

// LoadConfig loads configuration from file and environment variables.
// Priority: env vars > new config (~/.config/sm/) > legacy config (~/.config/supermemory/).
func LoadConfig() *Config {
	cfg := &Config{BaseURL: DefaultBaseURL}

	// Try legacy config first (~/.config/supermemory/config.json)
	home, _ := os.UserHomeDir()
	legacyPath := filepath.Join(home, ".config", "supermemory", "config.json")
	if data, err := os.ReadFile(legacyPath); err == nil {
		var legacy struct {
			APIKey  string `json:"apiKey"`
			APIBase string `json:"apiBase"`
		}
		if json.Unmarshal(data, &legacy) == nil {
			if legacy.APIKey != "" {
				cfg.APIKey = legacy.APIKey
			}
			if legacy.APIBase != "" {
				cfg.BaseURL = legacy.APIBase
			}
		}
	}

	// New config overrides legacy
	if data, err := os.ReadFile(ConfigPath()); err == nil {
		json.Unmarshal(data, cfg)
	}

	// Environment variables override everything
	if key := os.Getenv("SUPERMEMORY_API_KEY"); key != "" {
		cfg.APIKey = key
	}
	if base := os.Getenv("SUPERMEMORY_API_BASE"); base != "" {
		cfg.BaseURL = base
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}

	return cfg
}

// SaveConfig writes the config to ~/.config/sm/config.json.
func SaveConfig(cfg *Config) error {
	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	return os.WriteFile(path, append(data, '\n'), 0600)
}

// ReadConfigFile reads only the config file (no env override), for display purposes.
func ReadConfigFile() *Config {
	cfg := &Config{}
	if data, err := os.ReadFile(ConfigPath()); err == nil {
		json.Unmarshal(data, cfg)
	}
	return cfg
}
