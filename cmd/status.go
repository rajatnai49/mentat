package cmd

import (
	"github.com/fatih/color"
	"github.com/rajatnai49/mentat/parsers"
	"github.com/rajatnai49/mentat/ui"
	"github.com/rajatnai49/mentat/vault"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Aliases: []string{"st"},
	Short: "Personal task and knowledge management tool.",
	Run: func(cmd *cobra.Command, args []string) {
		color.Green("Status")
		nts := parsers.FolderIterator(cfg.VaultPath, parsers.DailyFileParser)
		generateSummary(nts)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func generateSummary(noteTasks []*vault.NoteTask) {
	var items []vault.TaskItem
	for _, nt := range noteTasks {
		for _, t := range nt.Tasks {
			if !t.Done {
				items = append(items, vault.TaskItem{
					Task:     t,
					Filepath: nt.FilePath,
				})
			}
		}
	}

	ui.RenderTasks(items)
}
