package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franksde/supermemory-cli/internal"
)

func TestDocGetEscapesDocumentIDPathSegment(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"a/b"}`))
	}))
	defer server.Close()

	writeTestConfig(t, server.URL, &internal.Config{})

	if err := docGetCmd.RunE(docGetCmd, []string{"a/b"}); err != nil {
		t.Fatalf("doc get: %v", err)
	}
	if gotPath != "/v3/documents/a%2Fb" {
		t.Fatalf("doc get path = %q, want /v3/documents/a%%2Fb", gotPath)
	}
}
