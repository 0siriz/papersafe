package shards

import (
	"fmt"
	"io"
	"strings"

	"github.com/0siriz/papersafe/internal/encoding/zbase32"
	"github.com/0siriz/papersafe/pkg/keyshard"
)

const mnemonicColumns = 3

func printHeader(w io.Writer, quorum *keyshard.Quorum, opts generateOptions) {
	mode := "unsealed"
	if quorum.IsSealed() {
		mode = "sealed"
	}

	_, _ = fmt.Fprintf(w, "%-10s %s\n", "Quorum ID:", zbase32.StdEncoding.EncodeToString(quorum.ID))
	_, _ = fmt.Fprintf(w, "%-10s %d (threshold %d)\n", "Shards:", opts.parts, opts.threshold)
	_, _ = fmt.Fprintf(w, "%-10s %s\n", "Mode:", mode)
	_, _ = fmt.Fprintf(w, "%-10s %s\n", "Output:", opts.outDir)
	_, _ = fmt.Fprintln(w)
}

func printShardMnemonic(w io.Writer, shardSet keyshard.ShardSet, path string, index, total int) {
	_, _ = fmt.Fprintf(w, "Shard %d of %d (ID %d) -> %s\n", index, total, shardSet.Shard.ID, path)
	writeMnemonic(w, shardSet.Words)
	_, _ = fmt.Fprintln(w)
}

func writeMnemonic(w io.Writer, words []string) {
	width := 0
	for _, word := range words {
		width = max(width, len(word))
	}

	for i := 0; i < len(words); i += mnemonicColumns {
		end := min(i+mnemonicColumns, len(words))

		columns := make([]string, 0, end-i)
		for j := i; j < end; j++ {
			columns = append(columns, fmt.Sprintf("%2d. %-*s", j+1, width, words[j]))
		}

		line := strings.TrimRight(strings.Join(columns, "   "), " ")
		_, _ = fmt.Fprintf(w, "  %s\n", line)
	}
}

func printSummary(w io.Writer, count int, outDir string) {
	noun := "shard PDFs"
	if count == 1 {
		noun = "shard PDF"
	}

	_, _ = fmt.Fprintf(w, "Wrote %d %s to %s.\n", count, noun, outDir)
	_, _ = fmt.Fprintln(w, "Write each mnemonic on its matching shard PDF, then clear your terminal.")
	_, _ = fmt.Fprintln(w, "Never store enough shards together to meet the reconstruction threshold.")
}
