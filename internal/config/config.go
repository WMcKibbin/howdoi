package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the howdoi configuration.
type Config struct {
	DefaultProvider string          `yaml:"default_provider"`
	Providers       ProvidersConfig `yaml:"providers"`
}

// ProvidersConfig holds per-provider settings.
type ProvidersConfig struct {
	GitHub GitHubConfig `yaml:"github"`
	Ollama OllamaConfig `yaml:"ollama"`
	Claude ClaudeConfig `yaml:"claude"`
}

// GitHubConfig holds GitHub Models provider settings.
type GitHubConfig struct {
	Token string `yaml:"token"`
	Model string `yaml:"model"`
}

// OllamaConfig holds Ollama provider settings.
type OllamaConfig struct {
	Host  string `yaml:"host"`
	Model string `yaml:"model"`
}

// ClaudeConfig holds Claude provider settings.
type ClaudeConfig struct {
	Model string `yaml:"model"`
}

// Path returns the default config file path.
func Path() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "howdoi", "config.yaml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "howdoi", "config.yaml")
}

// Load reads the config from disk, returning defaults if not found.
func Load() (*Config, error) {
	cfg := &Config{
		DefaultProvider: "claude",
	}

	data, err := os.ReadFile(Path())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the config to disk.
func Save(cfg *Config) error {
	p := Path()
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}
