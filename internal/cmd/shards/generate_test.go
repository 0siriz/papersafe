package shards

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateOptionsValidate(t *testing.T) {
	tests := []struct {
		name    string
		opts    generateOptions
		wantErr string
	}{
		{
			name: "valid",
			opts: generateOptions{parts: 5, threshold: 3, outDir: "."},
		},
		{
			name: "valid minimum",
			opts: generateOptions{parts: 2, threshold: 2, outDir: "."},
		},
		{
			name:    "threshold too low",
			opts:    generateOptions{parts: 5, threshold: 1},
			wantErr: "--threshold",
		},
		{
			name:    "parts too low",
			opts:    generateOptions{parts: 1, threshold: 2},
			wantErr: "--parts",
		},
		{
			name:    "parts below threshold",
			opts:    generateOptions{parts: 2, threshold: 3},
			wantErr: "cannot be less",
		},
		{
			name:    "too many parts",
			opts:    generateOptions{parts: 256, threshold: 2},
			wantErr: "255",
		},
		{
			name:    "empty output directory",
			opts:    generateOptions{parts: 2, threshold: 2, outDir: " "},
			wantErr: "--out-dir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.validate()

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, err)
			}
		})
	}
}

func TestEnsureTargetsAvailable(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "keyshard-1.pdf")

	if err := os.WriteFile(existing, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := ensureTargetsAvailable([]string{existing}, false)
	if err == nil {
		t.Fatal("expected an error for an existing target")
	}
	if !strings.Contains(err.Error(), "--force") {
		t.Fatalf("expected error to mention --force, got %q", err)
	}

	if err := ensureTargetsAvailable([]string{existing}, true); err != nil {
		t.Fatalf("expected --force to allow existing targets: %v", err)
	}

	missing := filepath.Join(dir, "keyshard-2.pdf")
	if err := ensureTargetsAvailable([]string{missing}, false); err != nil {
		t.Fatalf("expected missing target to be available: %v", err)
	}
}

func TestGenerateCommand(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "nested", "output")

	cmd := generateCommand()

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--parts", "2", "--threshold", "2", "--sealed", "--out-dir", outDir})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("generate command failed: %v", err)
	}

	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 shard PDFs, got %d", len(entries))
	}

	output := out.String()
	if !strings.Contains(output, "Mode:") || !strings.Contains(output, "sealed") {
		t.Errorf("output does not report sealed mode:\n%s", output)
	}
	if !strings.Contains(output, "Wrote 2 shard PDFs") {
		t.Errorf("output does not contain summary:\n%s", output)
	}
	if strings.Count(output, "Shard ") != 2 {
		t.Errorf("expected 2 shard blocks, got %d:\n%s", strings.Count(output, "Shard "), output)
	}
	if !strings.Contains(output, " 1. ") || !strings.Contains(output, "24. ") {
		t.Errorf("output does not contain a full numbered mnemonic:\n%s", output)
	}
}
