package cmd

import (
	"github.com/fatih/color"
	"github.com/rajatnai49/mentat/parsers"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:     "clean",
	Aliases: []string{"cln"},
	Short:   "Rename daily notes that have no pending tasks.",
	Long: `Rename daily notes that have no pending tasks.

Mentat scans daily markdown files named YYYYMMDD.md in your vault. If a file has
no unchecked checkbox tasks, it is renamed with an -X suffix, for example
20260516.md becomes 20260516-X.md.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		numOfFileChanges := 0
		err := parsers.DailyFilesIterator(cfg.VaultPath, func(path string) error {
			isRenamed, err := parsers.DailyFileCleaner(path)
			if err != nil {
				return err
			}
			if isRenamed {
				numOfFileChanges++
			}
			return nil
		})
		color.Green("Number of files changed: %v\n", numOfFileChanges)
		return err
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
