package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/rajatnai49/mentat/helpers"
	"github.com/spf13/cobra"
)

var (
	month    bool
	year     bool
	day      string
	template string
)

var dlCmd = &cobra.Command{
	Use:     "daily-note",
	Aliases: []string{"dl"},
	Short:   "Create or open a daily, monthly, or yearly note.",
	Long: `Create or open a daily, monthly, or yearly note.

	Mentat creates the note file in your configured vault if it does not already
	exist, then opens it in nvim. Use --day to choose a date, --month for a monthly
	note, or --year for a yearly note.

	Use -t <template-name> for creating note from your templates.

	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now()
		var t time.Time
		var err error

		if day != "" {
			t, err = helpers.ParseDate(day)
			if err != nil {
				return fmt.Errorf("Not valid date provided.")
			}
		} else {
			t = now
		}

		err = createOrOpenFile(t, template)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	dlCmd.Flags().BoolVarP(&month, "month", "m", false, "month files")
	dlCmd.Flags().BoolVarP(&year, "year", "y", false, "year files")
	dlCmd.Flags().StringVarP(&day, "day", "d", "", "daily files")
	dlCmd.Flags().StringVarP(&template, "template", "t", "", "template name")

	rootCmd.AddCommand(dlCmd)
}

func createOrOpenFile(t time.Time, template_name string) error {
	filename := t.Format("20060102")
	if month {
		filename = t.Format("200601")
	} else if year {
		filename = t.Format("2006")
	}

	path := cfg.VaultPath + "/" + filename + ".md"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if template_name != "" {
			template_path := cfg.VaultPath + "/templates/" + template_name + ".md"
			err := copyFile(template_path, path)
			if err != nil {
				return err
			}
		} else if err = os.WriteFile(path, []byte(""), 0644); err != nil {
			return fmt.Errorf("Error in file creation: %w", err)
		}
	}

	cfg, err := Load()
	if err != nil {
		return fmt.Errorf("Error in loading config: %w", err)
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
		return fmt.Errorf("Error in opening file: %w", err)
	}
	return nil
}

func copyFile(src_path, dest_path string) error {
	source, err := os.Open(src_path)
	if err != nil {
		return fmt.Errorf("Template file does not exist: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dest_path)
	if err != nil {
		return fmt.Errorf("Error in file creation: %w", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("Error in creating note from template: %w", err)
	}

	return nil
}
