package cmd

import (
	"fmt"
	"os"

	"github.com/WMcKibbin/howdoi/internal/config"
	"github.com/WMcKibbin/howdoi/internal/provider"
	"github.com/spf13/cobra"
)

var (
	flagProvider string
	flagModel    string
)

var rootCmd = &cobra.Command{
	Use:   "howdoi",
	Short: "AI-powered command-line assistant",
	Long:  "howdoi helps you find and understand shell commands using AI providers like Claude, GitHub Models, and Ollama.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagProvider, "provider", "", "AI provider to use (claude, github, ollama)")
	rootCmd.PersistentFlags().StringVar(&flagModel, "model", "", "Model to use with the provider")
}

// getProvider resolves the provider from flags and config.
func getProvider() (provider.Provider, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	initProviders(cfg)

	name := flagProvider
	if name == "" {
		name = cfg.DefaultProvider
	}
	if name == "" {
		name = "claude"
	}

	return provider.Get(name)
}

func initProviders(cfg *config.Config) {
	claudeModel := cfg.Providers.Claude.Model
	if flagModel != "" && (flagProvider == "" || flagProvider == "claude") {
		claudeModel = flagModel
	}
	provider.Register(&provider.ClaudeProvider{Model: claudeModel})

	githubToken := cfg.Providers.GitHub.Token
	if githubToken == "" {
		githubToken = os.Getenv("GITHUB_TOKEN")
	}
	githubModel := cfg.Providers.GitHub.Model
	if flagModel != "" && flagProvider == "github" {
		githubModel = flagModel
	}
	provider.Register(&provider.GitHubProvider{Token: githubToken, Model: githubModel})

	ollamaHost := cfg.Providers.Ollama.Host
	ollamaModel := cfg.Providers.Ollama.Model
	if flagModel != "" && flagProvider == "ollama" {
		ollamaModel = flagModel
	}
	provider.Register(&provider.OllamaProvider{Host: ollamaHost, Model: ollamaModel})
}
