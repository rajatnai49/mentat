package config

import (
	"os"

	"github.com/fatih/color"
)

type Config struct {
	VaultPath string
}

func Load() *Config {
	vault := mustEnv("MENTAT_VAULT")

	return &Config{
		VaultPath: vault,
	}
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		color.Red("Not found variable %s", key)
		os.Exit(1)
	}
	return val
}
