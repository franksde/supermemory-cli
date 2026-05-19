package cmd

import "testing"

func useTempHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	return home
}

func TestMaskKeyHandlesShortKeys(t *testing.T) {
	tests := map[string]string{
		"":          "(not set)",
		"x":         "***",
		"xy":        "xy***",
		"12345678":  "12***",
		"123456789": "1234***6789",
	}

	for key, want := range tests {
		got := maskKey(key)
		if got != want {
			t.Fatalf("maskKey(%q) = %q, want %q", key, got, want)
		}
	}
}

func TestConfigSettersRejectInvalidRanges(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T) error
	}{
		{
			name: "default limit zero",
			run: func(t *testing.T) error {
				useTempHome(t)
				return configSetDefaultLimitCmd.RunE(configSetDefaultLimitCmd, []string{"0"})
			},
		},
		{
			name: "score threshold below range",
			run: func(t *testing.T) error {
				useTempHome(t)
				return configSetScoreThresholdCmd.RunE(configSetScoreThresholdCmd, []string{"-0.1"})
			},
		},
		{
			name: "score threshold above range",
			run: func(t *testing.T) error {
				useTempHome(t)
				return configSetScoreThresholdCmd.RunE(configSetScoreThresholdCmd, []string{"1.1"})
			},
		},
		{
			name: "max content length negative",
			run: func(t *testing.T) error {
				useTempHome(t)
				return configSetMaxContentLengthCmd.RunE(configSetMaxContentLengthCmd, []string{"-1"})
			},
		},
		{
			name: "api timeout zero",
			run: func(t *testing.T) error {
				useTempHome(t)
				return configSetApiTimeoutCmd.RunE(configSetApiTimeoutCmd, []string{"0"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.run(t); err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}
