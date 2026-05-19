package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Document operations (add, batch, get, delete)",
}

var docAddCmd = &cobra.Command{
	Use:   "add <file|->",
	Short: "Add a document",
	Long: `Add a document from a file or stdin to Supermemory (via /v3/documents).

Examples:
  sm doc add README.md --containerTag my-project
  cat notes.txt | sm doc add - --customId my-notes`,
	Args: cobra.ExactArgs(1),
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
			content = string(data)
		} else {
			data, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			content = string(data)
		}

		payload := map[string]interface{}{"content": content}
		if ct := getContainerTag(cmd); ct != "" {
			payload["containerTag"] = ct
		}
		if cid, _ := cmd.Flags().GetString("customId"); cid != "" {
			payload["customId"] = cid
		}

		resp, err := c.Post("/v3/documents", payload)
		if err != nil {
			return err
		}
		fmt.Println(string(resp))
		return nil
	},
}

var docBatchCmd = &cobra.Command{
	Use:   "batch <manifest.json>",
	Short: "Batch add documents",
	Long: `Batch add documents from a JSON manifest file (via /v3/documents/batch).
Automatically chunks into batches of 600.

Examples:
  sm doc batch docs.json
  sm doc batch docs.json --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		c, err := newClient()
		if err != nil {
			return err
		}

		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("failed to open manifest: %w", err)
		}
		defer f.Close()

		var docs []map[string]interface{}
		if err := json.NewDecoder(f).Decode(&docs); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		const maxBatch = 600
		total := len(docs)
		start := 0
		for start < total {
			end := start + maxBatch
			if end > total {
				end = total
			}
			batch := docs[start:end]
			if dryRun {
				fmt.Printf("Would send batch %d-%d of %d (size %d)\n", start+1, end, total, len(batch))
			} else {
				payload := map[string]interface{}{"documents": batch}
				resp, err := c.Post("/v3/documents/batch", payload)
				if err != nil {
					return fmt.Errorf("batch %d-%d failed: %w", start+1, end, err)
				}
				fmt.Printf("Batch %d-%d: %s\n", start+1, end, resp)
			}
			start = end
			if !dryRun && end < total {
				time.Sleep(200 * time.Millisecond)
			}
		}
		return nil
	},
}

var docGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a document by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := c.Get("/v3/documents/" + url.PathEscape(args[0]))
		if err != nil {
			return err
		}
		var out interface{}
		if err := json.Unmarshal(resp, &out); err != nil {
			return fmt.Errorf("invalid JSON response: %w", err)
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	},
}

var docDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a document by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		_, err = c.Delete("/v3/documents/" + url.PathEscape(args[0]))
		if err != nil {
			return err
		}
		fmt.Println(`{"deleted": true}`)
		return nil
	},
}

func init() {
	docAddCmd.Flags().String("containerTag", "", "Container tag")
	docAddCmd.Flags().String("customId", "", "Custom ID for dedup")
	docBatchCmd.Flags().Bool("dry-run", false, "Print batch plan without sending")

	docCmd.AddCommand(docAddCmd)
	docCmd.AddCommand(docBatchCmd)
	docCmd.AddCommand(docGetCmd)
	docCmd.AddCommand(docDeleteCmd)
}
