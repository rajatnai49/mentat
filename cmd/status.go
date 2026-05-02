package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/rajatnai49/mentat/config"
	"github.com/rajatnai49/mentat/parsers"
	"github.com/rajatnai49/mentat/vault"
	"github.com/spf13/cobra"
)

var status = &cobra.Command{
	Use:   "status",
	Short: "Personal task and knowledge management tool.",
	Run: func(cmd *cobra.Command, args []string) {
		color.Green("Status")
		cfg = config.Load()
		nts := parsers.FolderIterator(cfg.VaultPath, parsers.DailyFileParser)
		generateSummary(cfg, nts)
	},
	Version: "1.0.0",
}

func init() {
	rootCmd.AddCommand(status)
}

type TaskItem struct {
	Task     vault.Task
	Filepath string
}

func generateSummary(cfg *config.Config, noteTasks []*vault.NoteTask) {
	var items []TaskItem
	for _, nt := range noteTasks {
		for _, t := range nt.Tasks {
			if !t.Done {
				items = append(items, TaskItem{
					Task:     t,
					Filepath: nt.FilePath,
				})
			}
		}
	}

	for i, itm := range items {
		color.Green("%d. %s", i+1, itm.Task.Title)
	}

	var choice int
	fmt.Println("> Go to: ")
	fmt.Scanln(&choice)

	if choice < 0 || choice > len(items) {
		color.Red("Invalid Selection")
		return
	}

	selected := items[choice-1]

	cmd := exec.Command("nvim", cfg.VaultPath+"/"+selected.Filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		color.Red("%v", err)
		return
	}
}
