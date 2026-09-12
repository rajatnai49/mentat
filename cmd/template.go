package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

var (
	template_name string
)

var templateCmd = &cobra.Command{
	Use:     "create-template",
	Aliases: []string{"tmpl"},
	Short:   "Create or open a template note.",
	Long: `Create or open a template note.

	Mentat create the note file in your configured vault '/templates' folder if it does not already exist, then opens it in configured editor.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		validName := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
		if !validName.MatchString(template_name) {
			return fmt.Errorf("Invalid filename: %q", template_name)
		}

		path := cfg.VaultPath + "/templates/" + template_name + ".md"

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("Error in template folder creation: %w", err)
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err = os.WriteFile(path, []byte(""), 0644); err != nil {
				return fmt.Errorf("Error in file creation: %w", err)
			}
		}

		cfg, err := Load()
		if err != nil {
			return fmt.Errorf("Error in config load: %w", err)
		}

		if cfg.Editor == "" {
			cfg.Editor = "vim"
		}

		openCmd := exec.Command(cfg.Editor, path)

		openCmd.Stdin = os.Stdin
		openCmd.Stdout = os.Stdout
		openCmd.Stderr = os.Stderr

		err = openCmd.Run()
		if err != nil {
			return fmt.Errorf("Error in file opening: %w", err)
		}
		return nil
	},
}

func init() {
	templateCmd.Flags().StringVarP(&template_name, "name", "n", "", "template file name")
	rootCmd.AddCommand(templateCmd)
}
