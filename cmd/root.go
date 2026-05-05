package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/rajatnai49/mentat/config"
	"github.com/spf13/cobra"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:     "mentant",
	Aliases: []string{"mnt"},
	Short:   "Personal task and knowledge management tool.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cfg == nil {
			cfg = config.Load()
		}
	},
	Version: "1.0.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("Fail to execute the command")
		os.Exit(1)
	}
}
