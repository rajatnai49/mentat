package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	VaultPath string
}

func Load() (*Config, error) {
	keys := []string{
		"MENTAT_VAULT",
	}

	values := make(map[string]string)
	var errs []error

	for _, k := range keys {
		val, err := getEnv(k)
		if err != nil {
			errs = append(errs)
			continue
		}
		values[k] = val
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &Config{
		VaultPath: values["MENTAT_VAULT"],
	}, nil
}

func getEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("Missing ENV variable: %s", key)
	}
	return val, nil
}
