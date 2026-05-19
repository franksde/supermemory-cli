package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/franksde/supermemory-cli/internal"
)

func writeTestConfig(t *testing.T, baseURL string, cfg *internal.Config) {
	t.Helper()
	home := useTempHome(t)
	cfg.APIKey = "sm_test_key"
	cfg.BaseURL = baseURL
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	path := filepath.Join(home, ".config", "sm", "config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("create config dir: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func TestSearchDocsUsesV3EvenWhenDefaultV4Enabled(t *testing.T) {
	defaultV4 := true
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[]}`))
	}))
	defer server.Close()

	writeTestConfig(t, server.URL, &internal.Config{DefaultV4: &defaultV4})

	if err := runSearch(searchDocsCmd, []string{"hello"}); err != nil {
		t.Fatalf("runSearch: %v", err)
	}
	if gotPath != "/v3/search" {
		t.Fatalf("search docs called %q, want /v3/search", gotPath)
	}
}

func TestRunSearchReturnsInvalidJSONError(t *testing.T) {
	defaultV4 := true
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	writeTestConfig(t, server.URL, &internal.Config{DefaultV4: &defaultV4})

	if err := runSearch(searchCmd, []string{"hello"}); err == nil {
		t.Fatalf("expected invalid JSON error")
	}
}
