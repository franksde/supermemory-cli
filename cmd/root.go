package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/franksde/supermemory-cli/client"
	"github.com/franksde/supermemory-cli/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sm",
	Short: "Supermemory CLI — memory for AI agents",
	Long:  "sm is a blazing-fast CLI for the Supermemory API, designed for AI agent integration.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(forgetCmd)
	rootCmd.AddCommand(docCmd)
	rootCmd.AddCommand(convCmd)
	rootCmd.AddCommand(containerCmd)
	rootCmd.AddCommand(configCmd)
}

// newClient creates a new API client from config, returning error if API key is missing.
func newClient() (*client.Client, error) {
	cfg := internal.LoadConfig()
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key not configured.\n  Run: sm config set-key <key>\n  Or set SUPERMEMORY_API_KEY environment variable")
	}
	c := client.New(cfg.BaseURL, cfg.APIKey)
	c.HTTPClient.Timeout = time.Duration(cfg.GetAPITimeout()) * time.Second
	return c, nil
}

// getContainerTag returns the container tag from flag or config default.
func getContainerTag(cmd *cobra.Command) string {
	tag, _ := cmd.Flags().GetString("containerTag")
	if tag != "" {
		return tag
	}
	cfg := internal.LoadConfig()
	return cfg.DefaultContainerTag
}
