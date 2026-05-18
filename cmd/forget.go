package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var forgetCmd = &cobra.Command{
	Use:   "forget <id-or-content>",
	Short: "Forget a memory",
	Long: `Forget (soft delete) a memory by ID or exact content match (via DELETE /v4/memories).
If the argument starts with "mem_", it is treated as a memory ID; otherwise as content.

Examples:
  sm forget mem_abc123 --containerTag my-project
  sm forget "outdated fact about the project" --reason "no longer relevant"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		containerTag := getContainerTag(cmd)
		payload := map[string]interface{}{}
		if containerTag != "" {
			payload["containerTag"] = containerTag
		}

		// If it looks like a memory ID (starts with "mem_"), use id field
		if len(args[0]) > 4 && args[0][:4] == "mem_" {
			payload["id"] = args[0]
		} else {
			payload["content"] = args[0]
		}

		if reason, _ := cmd.Flags().GetString("reason"); reason != "" {
			payload["reason"] = reason
		}

		resp, err := c.DeleteWithBody("/v4/memories", payload)
		if err != nil {
			return err
		}
		fmt.Println(string(resp))
		return nil
	},
}

func init() {
	forgetCmd.Flags().String("containerTag", "", "Container tag")
	forgetCmd.Flags().String("reason", "", "Reason for forgetting")
}
