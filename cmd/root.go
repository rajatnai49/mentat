package cmd

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rajatnai49/mentat/vault"
	"github.com/spf13/cobra"
)

var cfg *vault.Config
var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "mentat",
	Aliases: []string{"mnt"},
	Short:   "Manage daily notes and pending tasks from a local markdown vault.",
	Long: `Mentat is a terminal tool for working with a local markdown vault.

It can create or open daily, monthly, and yearly notes, scan daily notes for
pending checkbox tasks, and show those tasks in an interactive list.`,
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
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("Fail to execute the command with error: %v", err)
		os.Exit(1)
	}
}
