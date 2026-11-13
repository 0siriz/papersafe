package shards

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/0siriz/papersafe/pkg/keyshard"
)

func TestWriteMnemonic(t *testing.T) {
	words := make([]string, 24)
	for i := range words {
		words[i] = fmt.Sprintf("word%02d", i+1)
	}

	var out bytes.Buffer
	writeMnemonic(&out, words)

	lines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
	if len(lines) != 8 {
		t.Fatalf("expected 8 rows, got %d:\n%s", len(lines), out.String())
	}

	for i, word := range words {
		row := lines[i/mnemonicColumns]
		label := fmt.Sprintf("%2d. %s", i+1, word)
		if !strings.Contains(row, label) {
			t.Errorf("row %d does not contain %q: %q", i/mnemonicColumns+1, label, row)
		}
	}

	if strings.Contains(lines[0], words[3]) {
		t.Errorf("first row contains the fourth word: %q", lines[0])
	}
}

func TestPrintHeaderMode(t *testing.T) {
	tests := []struct {
		name   string
		sealed bool
		want   string
	}{
		{name: "unsealed", want: "unsealed"},
		{name: "sealed", sealed: true, want: "sealed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quorum := &keyshard.Quorum{ID: make([]byte, 32)}
			if tt.sealed {
				quorum.SetSealed()
			}

			var out bytes.Buffer
			printHeader(&out, quorum, generateOptions{parts: 3, threshold: 2, outDir: "."})

			want := fmt.Sprintf("%-10s %s", "Mode:", tt.want)
			if !strings.Contains(out.String(), want) {
				t.Errorf("header does not contain %q:\n%s", want, out.String())
			}
		})
	}
}

func TestPrintShardMnemonic(t *testing.T) {
	words := make([]string, 24)
	for i := range words {
		words[i] = fmt.Sprintf("word%02d", i+1)
	}

	var out bytes.Buffer
	printShardMnemonic(&out, keyshard.ShardSet{
		Shard: keyshard.KeyShard{ID: 42},
		Words: words,
	}, "keyshard-42.pdf", 1, 3)

	output := out.String()
	for _, want := range []string{
		"Shard 1 of 3 (ID 42) -> keyshard-42.pdf",
		words[0],
		words[23],
	} {
		if !strings.Contains(output, want) {
			t.Errorf("output does not contain %q:\n%s", want, output)
		}
	}
}
