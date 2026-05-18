package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var convCmd = &cobra.Command{
	Use:   "conv",
	Short: "Conversation operations",
}

var convIngestCmd = &cobra.Command{
	Use:   "ingest <conversation.json>",
	Short: "Ingest or update a conversation",
	Long: `Ingest a conversation from a JSON file (via /v4/conversations).

Examples:
  sm conv ingest chat.json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()

		var payload map[string]interface{}
		if err := json.NewDecoder(f).Decode(&payload); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		resp, err := c.Post("/v4/conversations", payload)
		if err != nil {
			return err
		}
		fmt.Println(string(resp))
		return nil
	},
}

func init() {
	convCmd.AddCommand(convIngestCmd)
}
