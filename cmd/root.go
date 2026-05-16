package cmd

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rajatnai49/mentat/vault"
	"github.com/spf13/cobra"
)

var cfg *vault.Config

var rootCmd = &cobra.Command{
	Use:     "mentat",
	Aliases: []string{"mnt"},
	Short:   "Personal task and knowledge management tool.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if strings.HasPrefix(
			cmd.CommandPath(),
			"mentat config",
		) {
			return nil
		}

		var err error
		if cfg == nil {
			cfg, err = Load()
			if err != nil {
				return err
			}
		}
		return nil
	},
	Version: "1.0.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("Fail to execute the command with error: %v", err)
		os.Exit(1)
	}
}
