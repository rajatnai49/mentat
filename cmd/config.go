package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Config struct {
	VaultPath string `toml:"vault_path"`
	Editor    string `toml:"editor"`
}

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Handling configs for the Mentat",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config file for Mentat",
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
			if ans != "y" || ans != "yes" {
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

		cfg := Config{
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
	Short: "See the current configuration of the Mentat",
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

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(showConfigCommand)
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

func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	var cfg Config

	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
