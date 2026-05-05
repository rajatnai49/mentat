package cmd

import (
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	month bool
	year  bool
	day   string
)

var dlCmd = &cobra.Command{
	Use:     "daily-note",
	Aliases: []string{"dl"},
	Short:   "Create or open daily note",
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		var t time.Time

		if day != "" {
		} else {
			t = now
		}

		createOrOpenFile(t)
	},
}

func init() {
	dlCmd.Flags().BoolVarP(&month, "month", "m", false, "month files")
	dlCmd.Flags().BoolVarP(&year, "year", "y", false, "year files")
	dlCmd.Flags().StringVarP(&day, "day", "d", "", "daily files")

	rootCmd.AddCommand(dlCmd)
}

func createOrOpenFile(t time.Time) {
	filename := t.Format("20060102")
	if month {
		filename = t.Format("200601")
	} else if year {
		filename = t.Format("2006")
	}

	path := cfg.VaultPath + "/" + filename + ".md"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte(""), 0644)
	}

	cmd := exec.Command("nvim", path)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		color.Red("%v", err)
		return
	}
}
