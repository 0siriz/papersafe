package shards

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/0siriz/papersafe/internal/paper"
	"github.com/0siriz/papersafe/pkg/keyshard"
	"github.com/spf13/cobra"
)

type generateOptions struct {
	parts     int
	threshold int
	outDir    string
	force     bool
	sealed    bool
}

func generateCommand() *cobra.Command {
	var opts generateOptions

	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate keyshards",
		Aliases: []string{"g", "gen"},
		Long: `Generate encrypted keyshard PDFs and print each shard's 24-word mnemonic
to the terminal. Write each mnemonic on its matching shard PDF and keep shards
separate so that no single location contains enough shards to meet the
reconstruction threshold.

Use --sealed to exclude the quorum signing key from the split so that
reconstruction cannot issue new shards.`,
		Example: `  papersafe shards generate --parts 5 --threshold 3 --out-dir ./shards
  papersafe shards generate --parts 5 --threshold 3 --sealed`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(cmd, opts)
		},
	}

	cmd.Flags().IntVarP(&opts.parts, "parts", "n", 0, "Number of keyshard parts")
	cmd.Flags().IntVarP(&opts.threshold, "threshold", "k", 0, "Threshold for reconstruction")
	cmd.Flags().StringVarP(&opts.outDir, "out-dir", "o", ".", "Directory to write keyshard PDFs")
	cmd.Flags().BoolVarP(&opts.force, "force", "f", false, "Overwrite existing keyshard PDFs")
	cmd.Flags().BoolVarP(&opts.sealed, "sealed", "s", false, "Seal the quorum so reconstruction cannot issue new shards")
	_ = cmd.MarkFlagRequired("parts")
	_ = cmd.MarkFlagRequired("threshold")

	return cmd
}

func (o generateOptions) validate() error {
	if o.threshold < 2 {
		return errors.New("--threshold must be at least 2")
	}
	if o.parts < 2 {
		return errors.New("--parts must be at least 2")
	}
	if o.parts < o.threshold {
		return fmt.Errorf("--parts (%d) cannot be less than --threshold (%d)", o.parts, o.threshold)
	}
	if o.parts > 255 {
		return errors.New("--parts cannot exceed 255")
	}
	if strings.TrimSpace(o.outDir) == "" {
		return errors.New("--out-dir cannot be empty")
	}

	return nil
}

func runGenerate(cmd *cobra.Command, opts generateOptions) error {
	if err := opts.validate(); err != nil {
		return err
	}

	if err := os.MkdirAll(opts.outDir, 0o700); err != nil {
		return fmt.Errorf("create output directory %q: %w", opts.outDir, err)
	}

	quorum, err := keyshard.NewQuorum()
	if err != nil {
		return err
	}

	if opts.sealed {
		quorum.SetSealed()
	}

	shardSets, err := quorum.MakeKeyshards(opts.parts, opts.threshold)
	if err != nil {
		return err
	}

	paths := make([]string, len(shardSets))
	for i, shardSet := range shardSets {
		paths[i] = filepath.Join(opts.outDir, fmt.Sprintf("keyshard-%d.pdf", shardSet.Shard.ID))
	}

	if err := ensureTargetsAvailable(paths, opts.force); err != nil {
		return err
	}

	out := cmd.OutOrStdout()

	printHeader(out, quorum, opts)

	for i, shardSet := range shardSets {
		m, err := paper.GetKeyshard(shardSet.Shard)
		if err != nil {
			return fmt.Errorf("prepare keyshard %d: %w", shardSet.Shard.ID, err)
		}

		document, err := m.Generate()
		if err != nil {
			return fmt.Errorf("generate keyshard %d PDF: %w", shardSet.Shard.ID, err)
		}

		if err := document.Save(paths[i]); err != nil {
			return fmt.Errorf("save keyshard %d to %q: %w", shardSet.Shard.ID, paths[i], err)
		}

		printShardMnemonic(out, shardSet, paths[i], i+1, len(shardSets))
	}

	printSummary(out, len(shardSets), opts.outDir)

	return nil
}

func ensureTargetsAvailable(paths []string, force bool) error {
	if force {
		return nil
	}

	var existing []string

	for _, path := range paths {
		_, err := os.Stat(path)
		switch {
		case err == nil:
			existing = append(existing, path)
		case errors.Is(err, os.ErrNotExist):
			continue
		default:
			return fmt.Errorf("check %q: %w", path, err)
		}
	}

	if len(existing) > 0 {
		return fmt.Errorf(
			"refusing to overwrite existing keyshard PDF(s): %s (use --force to overwrite)",
			strings.Join(existing, ", "),
		)
	}

	return nil
}
