package shards

import "github.com/spf13/cobra"

// NewCommand returns the "shards" subcommand.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "shards",
		Short:   "Manage keyshards",
		Aliases: []string{"s"},
	}

	cmd.AddCommand(generateCommand())

	return cmd
}
