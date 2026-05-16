package vault

type Config struct {
	VaultPath string `toml:"vault_path"`
	Editor    string `toml:"editor"`
}
