package cmd

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/rajatnai49/mentat/helpers"
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
	Short:   "Create or open a daily, monthly, or yearly note.",
	Long: `Create or open a daily, monthly, or yearly note.

Mentat creates the note file in your configured vault if it does not already
exist, then opens it in nvim. Use --day to choose a date, --month for a monthly
note, or --year for a yearly note.`,
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		var t time.Time
		var err error

		if day != "" {
			t, err = helpers.ParseDate(day)
			if err != nil {
				log.Fatalln(err)
				color.Red("Not valid date provided.")
				return
			}
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
		if err = os.WriteFile(path, []byte(""), 0644); err != nil {
			color.Red("Error in file creation: %v", err)
			return
		}
	}

	cfg, err := Load()
	if err != nil {
		color.Red("%v", err)
	}

	if cfg.Editor == "" {
		cfg.Editor = "vim"
	}

	cmd := exec.Command(cfg.Editor, path)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		color.Red("%v", err)
		return
	}
}

