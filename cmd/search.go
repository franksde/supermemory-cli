package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/franksde/supermemory-cli/internal"
	"github.com/spf13/cobra"
)

func truncateString(s string, maxLen int) string {
	if maxLen > 0 && len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

func filterResults(out interface{}, scoreThreshold float64, maxLen int) {
	m, ok := out.(map[string]interface{})
	if !ok {
		return
	}
	resultsVal, ok := m["results"].([]interface{})
	if !ok {
		return
	}

	var filtered []interface{}
	for _, item := range resultsVal {
		res, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		
		var score float64
		if simVal, ok := res["similarity"].(float64); ok {
			score = simVal
		} else if scoreVal, ok := res["score"].(float64); ok {
			score = scoreVal
		}

		if score < scoreThreshold {
			continue
		}

		if maxLen > 0 {
			if memVal, ok := res["memory"].(string); ok {
				res["memory"] = truncateString(memVal, maxLen)
			}
			if chunks, ok := res["chunks"].([]interface{}); ok {
				for _, ch := range chunks {
					if chunk, ok := ch.(map[string]interface{}); ok {
						if content, ok := chunk["content"].(string); ok {
							chunk["content"] = truncateString(content, maxLen)
						}
					}
				}
			}
		}

		filtered = append(filtered, res)
	}
	m["results"] = filtered
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search memories and documents",
	Long: `Search Supermemory using semantic search.
Defaults to v3 document search. Use --v4 for memory search (lower latency).

Examples:
  sm search "React authentication patterns"
  sm search "user preferences" --containerTag my-project --v4
  sm search docs "API rate limits" --limit 5`,
	Args: cobra.MinimumNArgs(1),
	RunE: runSearch,
}

var searchDocsCmd = &cobra.Command{
	Use:   "docs <query>",
	Short: "Search documents (backward-compatible alias)",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

func runSearch(cmd *cobra.Command, args []string) error {
	cfg := internal.LoadConfig()
	c, err := newClient()
	if err != nil {
		return err
	}

	query := strings.Join(args, " ")
	limit := cfg.GetDefaultLimit()
	if cmd.Flags().Changed("limit") {
		limit, _ = cmd.Flags().GetInt("limit")
	}

	containerTag := getContainerTag(cmd)
	chunkThresh, _ := cmd.Flags().GetFloat64("chunkThreshold")
	
	useV4 := cfg.GetDefaultV4()
	if cmd.Flags().Changed("v4") {
		useV4, _ = cmd.Flags().GetBool("v4")
	}

	payload := map[string]interface{}{"q": query, "limit": limit}

	if useV4 && containerTag != "" {
		payload["containerTag"] = containerTag
	} else if containerTag != "" {
		payload["containerTags"] = []string{containerTag}
	}
	if chunkThresh > 0 {
		payload["chunkThreshold"] = chunkThresh
	}

	endpoint := "/v3/search"
	if useV4 {
		endpoint = "/v4/search"
	}

	resp, err := c.Post(endpoint, payload)
	if err != nil {
		return err
	}

	var out interface{}
	json.Unmarshal(resp, &out)

	filterResults(out, cfg.GetScoreThreshold(), cfg.GetMaxContentLength())

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func init() {
	for _, cmd := range []*cobra.Command{searchCmd, searchDocsCmd} {
		cmd.Flags().Int("limit", 10, "Max results")
		cmd.Flags().Bool("v4", false, "Use v4 API (memory search, lower latency)")
		cmd.Flags().String("containerTag", "", "Container tag filter")
		cmd.Flags().Float64("chunkThreshold", 0, "Chunk selection sensitivity (0-1)")
	}
	searchCmd.AddCommand(searchDocsCmd)
}
