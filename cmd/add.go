package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <text|->",
	Short: "Add a memory",
	Long: `Add a memory directly to Supermemory (via /v4/memories).
Pass text as argument or '-' to read from stdin.

Examples:
  sm add "Frank prefers Claude for frontend work"
  sm add "Project uses React + Node.js" --containerTag my-project
  echo "memory from pipe" | sm add -`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		var content string
		if args[0] == "-" {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			content = strings.TrimSpace(string(data))
		} else {
			content = strings.Join(args, " ")
		}

		if content == "" {
			return fmt.Errorf("memory content cannot be empty")
		}

		containerTag := getContainerTag(cmd)
		if containerTag == "" {
			return fmt.Errorf("container tag required. Use --containerTag or set default: sm config set-tag <tag>")
		}

		memory := map[string]interface{}{"content": content}

		if tags, _ := cmd.Flags().GetString("tags"); tags != "" {
			memory["metadata"] = map[string]interface{}{
				"tags": tags,
			}
		}

		payload := map[string]interface{}{
			"memories":     []interface{}{memory},
			"containerTag": containerTag,
		}

		resp, err := c.Post("/v4/memories", payload)
		if err != nil {
			return err
		}
		fmt.Println(string(resp))
		return nil
	},
}

func init() {
	addCmd.Flags().String("containerTag", "", "Container tag for the memory")
	addCmd.Flags().String("tags", "", "Comma-separated tags (stored as metadata)")
}
