package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var editor_name string

var vaultCmd = &cobra.Command{
	Use: "vault",
	Aliases: []string{"vlt"},
	Short:   "Open mentat vault in the editor.",
	Long: `Open mentat vault in the editor.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := Load()
		if err != nil {
			return err
		}
		editor := cfg.Editor
		if editor_name != "" {
			editor = editor_name
		}

		op_cmd := exec.Command(editor, cfg.VaultPath)

		op_cmd.Stdin = os.Stdin
		op_cmd.Stdout = os.Stdout
		op_cmd.Stderr = os.Stderr

		err = op_cmd.Run()
		if err != nil {
			return fmt.Errorf("Error in opening vault: %w", err)
		}

		return nil
	},
}

func init() {
	vaultCmd.Flags().StringVarP(&editor_name, "editor", "e", "", "editor command")
	rootCmd.AddCommand(vaultCmd)
}
