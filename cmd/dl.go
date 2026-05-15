package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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
	Short:   "Create or open daily, month or year note",
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		var t time.Time
		var err error

		if day != "" {
			t, err = parseDate(day)
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

func parseDate(input string) (time.Time, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	t := time.Now()

	switch input {
	case "today":
		return t, nil
	case "yesterday":
		return t.AddDate(0, 0, -1), nil
	case "tomorrow":
		return t.AddDate(0, 0, 1), nil
	}

	allowedFormats := []string{
		"2006-01-02",
		"20060102",
		"02-01-2006",
		"02012006",
	}

	var err error
	for _, f := range allowedFormats {
		t, err = time.Parse(f, input)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("Invalid Date")
}
