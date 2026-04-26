package config

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

type Config struct {
	VaultPath string
	DBPath    string
}

func Load() *Config {
	vault := mustEnv("MENTAT_VAULT")

	db := os.Getenv("MENTAT_DB")

	if db == "" {
		db = filepath.Join(vault, "/.mentat.db")
	}

	return &Config{
		VaultPath: vault,
		DBPath:    db,
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
