package shards

import (
	"fmt"

	"github.com/0siriz/papersafe/internal/paper"
	"github.com/0siriz/papersafe/pkg/keyshard"
	"github.com/spf13/cobra"
)

func generateCommand() *cobra.Command {
	var parts int
	var threshold int

	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate keyshards",
		Aliases: []string{"g", "gen"},
		RunE: func(cmd *cobra.Command, args []string) error {
			quorum, err := keyshard.NewQuorum()
			if err != nil {
				return err
			}

			shardSets, err := quorum.MakeKeyshards(parts, threshold)
			if err != nil {
				return err
			}

			for _, shardSet := range shardSets {
				shard := shardSet.Shard
				m, err := paper.GetKeyshard(shard)
				if err != nil {
					return err
				}

				document, err := m.Generate()
				if err != nil {
					return err
				}

				if err := document.Save(fmt.Sprintf("keyshard-%d.pdf", shard.ID)); err != nil {
					return err
				}

			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&parts, "parts", "n", 0, "Number of keyshard parts")
	cmd.Flags().IntVarP(&threshold, "threshold", "k", 0, "Threshold for reconstruction")
	_ = cmd.MarkFlagRequired("parts")
	_ = cmd.MarkFlagRequired("threshold")

	return cmd
}
