package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveConfigUsesPrivatePermissions(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := SaveConfig(&Config{APIKey: "sm_test"}); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	dirInfo, err := os.Stat(filepath.Join(home, ".config", "sm"))
	if err != nil {
		t.Fatalf("stat config dir: %v", err)
	}
	if got := dirInfo.Mode().Perm(); got != 0o700 {
		t.Fatalf("config dir permissions = %o, want 700", got)
	}

	fileInfo, err := os.Stat(ConfigPath())
	if err != nil {
		t.Fatalf("stat config file: %v", err)
	}
	if got := fileInfo.Mode().Perm(); got != 0o600 {
		t.Fatalf("config file permissions = %o, want 600", got)
	}
}
