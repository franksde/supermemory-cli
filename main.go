package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, bytes.TrimSpace(body))
	}
	return body, nil
}

func (c *Client) Get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) Post(path string, payload interface{}) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) Patch(path string, payload interface{}) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PATCH", c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) Delete(path string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func loadConfig() (baseURL, apiKey string) {
	apiKey = os.Getenv("SUPERMEMORY_API_KEY")
	baseURL = os.Getenv("SUPERMEMORY_API_BASE")
	if baseURL == "" {
		baseURL = "https://api.supermemory.ai"
	}
	if apiKey == "" {
		cfgPath := os.ExpandEnv("$HOME/.config/supermemory/config.json")
		if data, err := os.ReadFile(cfgPath); err == nil {
			var cfg struct {
				APIKey  string `json:"apiKey"`
				APIBase string `json:"apiBase"`
			}
			if json.Unmarshal(data, &cfg) == nil {
				if apiKey == "" && cfg.APIKey != "" {
					apiKey = cfg.APIKey
				}
				if cfg.APIBase != "" {
					baseURL = cfg.APIBase
				}
			}
		}
	}
	return baseURL, apiKey
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "supermemory",
		Short: "Supermemory.ai CLI",
		Long:  "Command-line client for Supermemory API",
	}

	docCmd := &cobra.Command{Use: "doc", Short: "Document operations"}
	convCmd := &cobra.Command{Use: "conv", Short: "Conversation operations"}
	searchCmd := &cobra.Command{Use: "search", Short: "Search operations"}
	containerCmd := &cobra.Command{Use: "container", Short: "Container tag operations"}

	// doc add
	docAddCmd := &cobra.Command{
		Use:   "add <file|-",
		Short: "Add a document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL, apiKey := loadConfig()
			if apiKey == "" {
				return fmt.Errorf("API key not set (SUPERMEMORY_API_KEY or config file)")
			}
			client := NewClient(baseURL, apiKey)
			var content string
			if args[0] == "-" {
				data, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				content = string(data)
			} else {
				data, err := os.ReadFile(args[0])
				if err != nil {
					return err
				}
				content = string(data)
			}
			payload := map[string]interface{}{"content": content}
			if ct, _ := cmd.Flags().GetString("containerTag"); ct != "" {
				payload["containerTag"] = ct
			}
			if cid, _ := cmd.Flags().GetString("customId"); cid != "" {
				payload["customId"] = cid
			}
			resp, err := client.Post("/v3/documents", payload)
			if err != nil {
				return err
			}
			fmt.Println(string(resp))
			return nil
		},
	}
	docAddCmd.Flags().String("containerTag", "", "Container tag")
	docAddCmd.Flags().String("customId", "", "Custom ID for dedup")
	docCmd.AddCommand(docAddCmd)

	// doc batch
	docBatchCmd := &cobra.Command{
		Use:   "batch <manifest.json>",
		Short: "Batch add documents",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			baseURL, apiKey := loadConfig()
			if apiKey == "" {
				return fmt.Errorf("API key not set")
			}
			client := NewClient(baseURL, apiKey)
			f, err := os.Open(args[0])
			if err != nil {
				return err
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
					resp, err := client.Post("/v3/documents/batch", payload)
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
	docBatchCmd.Flags().Bool("dry-run", false, "Print batch plan without sending")
	docCmd.AddCommand(docBatchCmd)

	// conv ingest
	convIngestCmd := &cobra.Command{
		Use:   "ingest <conversation.json>",
		Short: "Ingest or update a conversation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL, apiKey := loadConfig()
			if apiKey == "" {
				return fmt.Errorf("API key not set")
			}
			client := NewClient(baseURL, apiKey)
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()
			var payload map[string]interface{}
			if err := json.NewDecoder(f).Decode(&payload); err != nil {
				return fmt.Errorf("invalid JSON: %w", err)
			}
			resp, err := client.Post("/v4/conversations", payload)
			if err != nil {
				return err
			}
			fmt.Println(string(resp))
			return nil
		},
	}
	convCmd.AddCommand(convIngestCmd)

	// search docs
	searchDocsCmd := &cobra.Command{
		Use:   "docs <query>",
		Short: "Search documents",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL, apiKey := loadConfig()
			if apiKey == "" {
				return fmt.Errorf("API key not set")
			}
			client := NewClient(baseURL, apiKey)
			query := strings.Join(args, " ")
			limit, _ := cmd.Flags().GetInt("limit")
			containerTag, _ := cmd.Flags().GetString("containerTag")
			chunkThresh, _ := cmd.Flags().GetFloat64("chunkThreshold")
			payload := map[string]interface{}{"q": query, "limit": limit}
			if containerTag != "" {
				payload["containerTags"] = []string{containerTag}
			}
			if chunkThresh > 0 {
				payload["chunkThreshold"] = chunkThresh
			}
			resp, err := client.Post("/v3/search", payload)
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
	searchDocsCmd.Flags().Int("limit", 10, "Max results")
	searchDocsCmd.Flags().String("containerTag", "", "Container tag filter")
	searchDocsCmd.Flags().Float64("chunkThreshold", 0, "Chunk selection sensitivity (0-1)")
	searchCmd.AddCommand(searchDocsCmd)

	// container commands
	containerGetCmd := &cobra.Command{
		Use:   "get <containerTag>",
		Short: "Get container tag settings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL, apiKey := loadConfig()
			if apiKey == "" {
				return fmt.Errorf("API key not set")
			}
			client := NewClient(baseURL, apiKey)
			resp, err := client.Get("/v3/container-tags/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println(string(resp))
			return nil
		},
	}

	containerSetCmd := &cobra.Command{
		Use:   "set <containerTag> --settings key=value",
		Short: "Update container tag settings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL, apiKey := loadConfig()
			if apiKey == "" {
				return fmt.Errorf("API key not set")
			}
			client := NewClient(baseURL, apiKey)
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
			resp, err := client.Patch("/v3/container-tags/"+args[0], payload)
			if err != nil {
				return err
			}
			fmt.Println(string(resp))
			return nil
		},
	}

	containerMergeCmd := &cobra.Command{
		Use:   "merge <sourceTag> <targetTag>",
		Short: "Merge container tags",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL, apiKey := loadConfig()
			if apiKey == "" {
				return fmt.Errorf("API key not set")
			}
			client := NewClient(baseURL, apiKey)
			payload := map[string]interface{}{
				"sourceTags": []string{args[0]},
				"targetTag":  args[1],
			}
			resp, err := client.Post("/v3/container-tags/merge", payload)
			if err != nil {
				return err
			}
			fmt.Println(string(resp))
			return nil
		},
	}

	containerDeleteCmd := &cobra.Command{
		Use:   "delete <containerTag>",
		Short: "Delete container tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL, apiKey := loadConfig()
			if apiKey == "" {
				return fmt.Errorf("API key not set")
			}
			client := NewClient(baseURL, apiKey)
			resp, err := client.Delete("/v3/container-tags/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println(string(resp))
			return nil
		},
	}

	containerCmd.AddCommand(containerGetCmd)
	containerCmd.AddCommand(containerSetCmd)
	containerCmd.AddCommand(containerMergeCmd)
	containerCmd.AddCommand(containerDeleteCmd)
	containerSetCmd.Flags().String("settings", "", "Settings as key=value pairs (comma-separated)")

	rootCmd.AddCommand(docCmd)
	rootCmd.AddCommand(convCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(containerCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}