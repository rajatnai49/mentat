package cmd

import (
	"github.com/rajatnai49/mentat/parsers"
	"github.com/rajatnai49/mentat/ui"
	"github.com/rajatnai49/mentat/vault"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st"},
	Short:   "Get status of all pending tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		ui.RenderList(getTaskItems)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func getTaskItems() []vault.TaskItem {
	var nts []*vault.NoteTask

	parsers.DailyFilesIterator(cfg.VaultPath, func (path string) error {
		nt, err := parsers.DailyFileParser(path)
		if err != nil {
			return err
		}

		if nt != nil {
			nts = append(nts, nt)
		}

		return nil
	})

	var items []vault.TaskItem
	for _, nt := range nts {
		for _, t := range nt.Tasks {
			if !t.Done {
				items = append(items, vault.TaskItem{
					Task:     t,
					Filepath: nt.FilePath,
				})
			}
		}
	}
	return items
}
