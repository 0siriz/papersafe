package main

import (
	"fmt"
	"os"

	"github.com/0siriz/papersafe/internal/cmd/backup"
	"github.com/0siriz/papersafe/internal/cmd/restore"
	"github.com/0siriz/papersafe/internal/cmd/shards"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "papersafe",
	Short:         "Papersafe: Generate secure offline paper backups",
	SilenceErrors: true,
	SilenceUsage:  true,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

// main is the papersafe CLI entry point.
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(backup.NewCommand())
	rootCmd.AddCommand(restore.NewCommand())
	rootCmd.AddCommand(shards.NewCommand())
}
