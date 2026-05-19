package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Container tag operations (get, set, merge, delete)",
}

var containerGetCmd = &cobra.Command{
	Use:   "get <containerTag>",
	Short: "Get container tag settings",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := c.Get("/v3/container-tags/" + url.PathEscape(args[0]))
		if err != nil {
			return err
		}
		fmt.Println(string(resp))
		return nil
	},
}

var containerSetCmd = &cobra.Command{
	Use:   "set <containerTag> --settings key=value,...",
	Short: "Update container tag settings",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		settingsStr, _ := cmd.Flags().GetString("settings")
		if settingsStr == "" {
			return fmt.Errorf("--settings required")
		}
		pairs := strings.Split(settingsStr, ",")
		settings := make(map[string]interface{})
		for _, p := range pairs {
			kv := strings.SplitN(p, "=", 2)
			if len(kv) != 2 {
				return fmt.Errorf("invalid setting: %s", p)
			}
			settings[kv[0]] = kv[1]
		}
		payload := map[string]interface{}{"settings": settings}
		resp, err := c.Patch("/v3/container-tags/"+url.PathEscape(args[0]), payload)
		if err != nil {
			return err
		}
		fmt.Println(string(resp))
		return nil
	},
}

var containerMergeCmd = &cobra.Command{
	Use:   "merge <sourceTag> <targetTag>",
	Short: "Merge container tags",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		payload := map[string]interface{}{
			"sourceTags": []string{args[0]},
			"targetTag":  args[1],
		}
		resp, err := c.Post("/v3/container-tags/merge", payload)
		if err != nil {
			return err
		}
		fmt.Println(string(resp))
		return nil
	},
}

var containerDeleteCmd = &cobra.Command{
	Use:   "delete <containerTag>",
	Short: "Delete container tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := c.Delete("/v3/container-tags/" + url.PathEscape(args[0]))
		if err != nil {
			return err
		}
		fmt.Println(string(resp))
		return nil
	},
}

func init() {
	containerSetCmd.Flags().String("settings", "", "Settings as key=value pairs (comma-separated)")

	containerCmd.AddCommand(containerGetCmd)
	containerCmd.AddCommand(containerSetCmd)
	containerCmd.AddCommand(containerMergeCmd)
	containerCmd.AddCommand(containerDeleteCmd)
}
