package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the CLI configuration.
type Config struct {
	APIKey              string   `json:"api_key,omitempty"`
	BaseURL             string   `json:"base_url,omitempty"`
	DefaultContainerTag string   `json:"default_container_tag,omitempty"`
	DefaultLimit        *int     `json:"default_limit,omitempty"`
	ScoreThreshold      *float64 `json:"score_threshold,omitempty"`
	MaxContentLength    *int     `json:"max_content_length,omitempty"`
	APITimeout          *int     `json:"api_timeout,omitempty"`
	DefaultV4           *bool    `json:"default_v4,omitempty"`
}

func (c *Config) GetDefaultLimit() int {
	if c.DefaultLimit != nil {
		return *c.DefaultLimit
	}
	return 3
}

func (c *Config) GetScoreThreshold() float64 {
	if c.ScoreThreshold != nil {
		return *c.ScoreThreshold
	}
	return 0.0
}

func (c *Config) GetMaxContentLength() int {
	if c.MaxContentLength != nil {
		return *c.MaxContentLength
	}
	return 800
}

func (c *Config) GetAPITimeout() int {
	if c.APITimeout != nil {
		return *c.APITimeout
	}
	return 5
}

func (c *Config) GetDefaultV4() bool {
	if c.DefaultV4 != nil {
		return *c.DefaultV4
	}
	return true
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
