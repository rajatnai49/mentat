package cmd

import (
	"github.com/fatih/color"
	"github.com/rajatnai49/mentat/parsers"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:     "clean",
	Aliases: []string{"cln"},
	Short:   "",
	Run: func(cmd *cobra.Command, args []string) {
		numOfFileChanges := 0
		parsers.DailyFilesIterator(cfg.VaultPath, func(path string) error {
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
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
