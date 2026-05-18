package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent memories",
	Long: `List memory entries from a container tag (via /v4/memories/list).
Requires --containerTag or a default container tag in config.

Examples:
  sm list --containerTag my-project
  sm list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		containerTag := getContainerTag(cmd)
		if containerTag == "" {
			return fmt.Errorf("container tag required. Use --containerTag or set default: sm config set-tag <tag>")
		}

		payload := map[string]interface{}{
			"containerTags": []string{containerTag},
		}

		resp, err := c.Post("/v4/memories/list", payload)
		if err != nil {
			return err
		}

		var out interface{}
		json.Unmarshal(resp, &out)
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	},
}

func init() {
	listCmd.Flags().String("containerTag", "", "Container tag to list memories from")
}
