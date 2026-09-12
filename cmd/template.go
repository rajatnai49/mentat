package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	list          bool
	template_name string
)

var templateCmd = &cobra.Command{
	Use:     "create-template",
	Aliases: []string{"tmpl"},
	Short:   "Create or open a note template.",
	Long: `Create or open a note template.

	Mentat stores templates in the templates/ folder of your configured vault. If the template does not exist, Mentat creates it, then opens it in your configured editor.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if list {
			return listTemplates()
		}

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

		openCmd := exec.Command(cfg.Editor, path)

		openCmd.Stdin = os.Stdin
		openCmd.Stdout = os.Stdout
		openCmd.Stderr = os.Stderr

		err := openCmd.Run()
		if err != nil {
			return fmt.Errorf("Error in file opening: %w", err)
		}
		return nil
	},
}

func listTemplates() error {
	template_folder := cfg.VaultPath + "/templates"
	entries, err := filepath.Glob(template_folder + "/*.md")
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		color.Red("No templates founds\n")
	} else {
		color.Green("Your templates:\n")
		total_templates := 0
		for _, v := range entries {
			name := filepath.Base(v)
			fmt.Println(name)
			total_templates++
		}
		color.Yellow("\nYou have total: %d templates", total_templates)
	}

	return nil
}

func init() {
	templateCmd.Flags().BoolVarP(&list, "list", "l", false, "list templates")
	templateCmd.Flags().StringVarP(&template_name, "name", "n", "", "template file name")
	rootCmd.AddCommand(templateCmd)
}
