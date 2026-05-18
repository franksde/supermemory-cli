package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

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
			"api_key":               maskKey(cfg.APIKey),
			"base_url":              cfg.BaseURL,
			"default_container_tag": cfg.DefaultContainerTag,
			"default_limit":         cfg.GetDefaultLimit(),
			"score_threshold":       cfg.GetScoreThreshold(),
			"max_content_length":    cfg.GetMaxContentLength(),
			"api_timeout":           cfg.GetAPITimeout(),
			"default_v4":            cfg.GetDefaultV4(),
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
	configCmd.AddCommand(configSetDefaultLimitCmd)
	configCmd.AddCommand(configSetScoreThresholdCmd)
	configCmd.AddCommand(configSetMaxContentLengthCmd)
	configCmd.AddCommand(configSetApiTimeoutCmd)
	configCmd.AddCommand(configSetDefaultV4Cmd)
	configCmd.AddCommand(configShowCmd)
}

var configSetDefaultLimitCmd = &cobra.Command{
	Use:   "set-default-limit <limit>",
	Short: "Set the default search return limit",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		val, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid integer for limit: %w", err)
		}
		cfg := internal.ReadConfigFile()
		cfg.DefaultLimit = &val
		if err := internal.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "default_limit set to %d\n", val)
		return nil
	},
}

var configSetScoreThresholdCmd = &cobra.Command{
	Use:   "set-score-threshold <threshold>",
	Short: "Set the relevance score threshold (0.0 to 1.0)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		val, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return fmt.Errorf("invalid float for threshold: %w", err)
		}
		cfg := internal.ReadConfigFile()
		cfg.ScoreThreshold = &val
		if err := internal.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "score_threshold set to %v\n", val)
		return nil
	},
}

var configSetMaxContentLengthCmd = &cobra.Command{
	Use:   "set-max-content-length <length>",
	Short: "Set the maximum character count for content snippets",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		val, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid integer for length: %w", err)
		}
		cfg := internal.ReadConfigFile()
		cfg.MaxContentLength = &val
		if err := internal.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "max_content_length set to %d\n", val)
		return nil
	},
}

var configSetApiTimeoutCmd = &cobra.Command{
	Use:   "set-api-timeout <seconds>",
	Short: "Set the API timeout in seconds",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		val, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid integer for timeout: %w", err)
		}
		cfg := internal.ReadConfigFile()
		cfg.APITimeout = &val
		if err := internal.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "api_timeout set to %d\n", val)
		return nil
	},
}

var configSetDefaultV4Cmd = &cobra.Command{
	Use:   "set-default-v4 <true|false>",
	Short: "Set whether to use v4 API by default",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		val, err := strconv.ParseBool(args[0])
		if err != nil {
			return fmt.Errorf("invalid boolean: %w", err)
		}
		cfg := internal.ReadConfigFile()
		cfg.DefaultV4 = &val
		if err := internal.SaveConfig(cfg); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "default_v4 set to %v\n", val)
		return nil
	},
}
