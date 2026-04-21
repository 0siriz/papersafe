package main

import (
	"github.com/0siriz/papersafe/internal/cmd/backup"
	"github.com/0siriz/papersafe/internal/cmd/restore"
	"github.com/0siriz/papersafe/internal/cmd/shards"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "papersafe",
	Short: "Papersafe: Generate secure offline paper backups",
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(backup.NewCommand())
	rootCmd.AddCommand(restore.NewCommand())
	rootCmd.AddCommand(shards.NewCommand())
}
