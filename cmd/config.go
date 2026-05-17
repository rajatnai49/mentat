package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/rajatnai49/mentat/vault"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Manage Mentat configuration.",
	Long: `Manage Mentat configuration.

The config file stores your markdown vault path and preferred editor. Mentat
uses this config for task scanning and opening files from interactive views.`,
}

var initConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "Create the Mentat config file.",
	Long: `Create the Mentat config file.

Mentat will ask for your vault path and preferred editor, then write the config
to your user config directory at mentat/config.toml.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			color.Red("Not be able to find config dir")
			return err
		}

		mentatDir := filepath.Join(userConfigDir, "mentat")
		err = os.MkdirAll(mentatDir, 0755)
		if err != nil {
			color.Red("Not be able to create mentat config folder")
			return err
		}

		mentatFile := filepath.Join(mentatDir, "config.toml")

		if _, err := os.Stat(mentatFile); err == nil {
			var ans string
			color.Yellow("Config already exist, %v", mentatFile)
			color.Yellow("Overwrite [y/n]: ")
			fmt.Scanln(&ans)
			ans = strings.ToLower(strings.TrimSpace(ans))
			if ans != "y" && ans != "yes" {
				color.Red("Abort mentat init")
				return nil
			}
		}

		var path, editor string

		color.Green("Vault path (required): ")
		fmt.Scanln(&path)
		path = strings.TrimSpace(path)
		if path == "" {
			color.Red("Vault path require for the config")
			return nil
		}

		color.Green("Editor you prefer [vim]: ")
		fmt.Scanln(&editor)
		editor = strings.TrimSpace(editor)

		if editor == "" {
			editor = "vim"
		}

		cfg := vault.Config{
			VaultPath: path,
			Editor:    editor,
		}

		var buf bytes.Buffer

		err = toml.NewEncoder(&buf).Encode(cfg)
		if err != nil {
			color.Red("Error in the saving file")
			return err
		}

		err = os.WriteFile(mentatFile, buf.Bytes(), 0644)
		if err != nil {
			color.Red("Error in the saving config file")
			return err
		}

		return nil
	},
}

var showConfigCommand = &cobra.Command{
	Use:   "show",
	Short: "Show the current Mentat configuration.",
	Long: `Show the current Mentat configuration.

Prints the config file path, vault path, and editor currently used by Mentat.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := ConfigPath()
		if err != nil {
			return err
		}

		color.Green("Config Path %s", configPath)

		if _, err = os.Stat(configPath); err != nil {
			if os.IsNotExist(err) {
				color.Red("Config file does not exist")
				return nil
			}

			return err
		}

		cfg, err := Load()
		if err != nil {
			return err
		}

		fmt.Println()

		color.Green("Vault Path %s", cfg.VaultPath)
		color.Green("Editor %s", cfg.Editor)

		return nil
	},
}

var openConfig = &cobra.Command{
	Use:   "open",
	Short: "Open the Mentat config file in your editor.",
	Long: `Open the Mentat config file in your configured editor.

If no editor is configured, Mentat falls back to vim.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := ConfigPath()
		if err != nil {
			return err
		}

		if _, err = os.Stat(configPath); err != nil {
			if os.IsNotExist(err) {
				color.Red("Config file does not exist")
				return nil
			}

			return err
		}

		cfg, err := Load()
		if err != nil {
			return err
		}

		if cfg.Editor == "" {
			cfg.Editor = "vim"
		}

		openCmd := exec.Command(cfg.Editor, configPath)

		openCmd.Stdin = os.Stdin
		openCmd.Stdout = os.Stdout
		openCmd.Stderr = os.Stderr

		err = openCmd.Run()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(initConfigCmd)
	configCmd.AddCommand(showConfigCommand)
	configCmd.AddCommand(openConfig)
	rootCmd.AddCommand(configCmd)
}

func ConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "Not be able to access Config dir", err
	}

	return filepath.Join(
		dir,
		"mentat",
		"config.toml",
	), nil
}

func Load() (*vault.Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	var cfg vault.Config

	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
