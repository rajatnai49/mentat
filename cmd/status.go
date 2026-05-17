package cmd

import (
	"github.com/fatih/color"
	"github.com/rajatnai49/mentat/parsers"
	"github.com/rajatnai49/mentat/ui"
	"github.com/rajatnai49/mentat/vault"
	"github.com/spf13/cobra"
)

var e bool

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st"},
	Short:   "Show all pending tasks from daily notes.",
	Long: `Show all pending tasks from daily notes.

Mentat scans markdown files named YYYYMMDD.md in your vault, extracts unchecked
checkbox tasks, and displays them in an interactive list. Select a task and
press enter to open its source note in your configured editor.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := ui.RenderList(cfg, getTaskItems)
		if e && err != nil {
			return err
		}
		return nil
	},
}

func init() {
	statusCmd.Flags().BoolVarP(&e, "error", "e", false, "Stop on error")
	rootCmd.AddCommand(statusCmd)
}

func getTaskItems() ([]vault.TaskItem, error) {
	var nts []*vault.NoteTask

	err := parsers.DailyFilesIterator(cfg.VaultPath, func(path string) error {
		nt, err := parsers.DailyFileParser(path)
		if err != nil {
			return err
		}

		if nt != nil {
			nts = append(nts, nt)
		}

		return nil
	})

	if e && err != nil {
		color.Red("%v", err)
		return nil, err
	}

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
	return items, nil
}
