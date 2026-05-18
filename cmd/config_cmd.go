package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/franksde/supermemory-cli/internal"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management (set-key, set-tag, show)",
}

var configSetKeyCmd = &cobra.Command{
	Use:   "set-key <api-key>",
	Short: "Set the Supermemory API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := internal.ReadConfigFile()
		cfg.APIKey = args[0]
		if cfg.BaseURL == "" {
			cfg.BaseURL = internal.DefaultBaseURL
		}
		if err := internal.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "API key saved to %s\n", internal.ConfigPath())
		return nil
	},
}

var configSetTagCmd = &cobra.Command{
	Use:   "set-tag <container-tag>",
	Short: "Set the default container tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := internal.ReadConfigFile()
		cfg.DefaultContainerTag = args[0]
		if cfg.BaseURL == "" {
			cfg.BaseURL = internal.DefaultBaseURL
		}
		if err := internal.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Default container tag set to %q\n", args[0])
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := internal.LoadConfig()
		// Mask the API key for display
		display := map[string]interface{}{
			"config_path":           internal.ConfigPath(),
			"api_key":              maskKey(cfg.APIKey),
			"base_url":             cfg.BaseURL,
			"default_container_tag": cfg.DefaultContainerTag,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(display)
	},
}

func maskKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 8 {
		return key[:2] + "***"
	}
	return key[:4] + "***" + key[len(key)-4:]
}

func init() {
	configCmd.AddCommand(configSetKeyCmd)
	configCmd.AddCommand(configSetTagCmd)
	configCmd.AddCommand(configShowCmd)
}
