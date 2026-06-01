package shards

import (
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

			// TODO: Make PDF files
			_ = shardSets

			return nil
		},
	}

	cmd.Flags().IntVarP(&parts, "parts", "n", 0, "Number of keyshard parts")
	cmd.Flags().IntVarP(&threshold, "threshold", "k", 0, "Threshold for reconstruction")
	_ = cmd.MarkFlagRequired("parts")
	_ = cmd.MarkFlagRequired("threshold")

	return cmd
}
