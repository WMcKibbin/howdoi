package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/WMcKibbin/howdoi/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure howdoi settings",
	Long:  "Interactive configuration wizard for howdoi settings.",
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Default provider [%s] (claude/github/ollama): ", cfg.DefaultProvider)
	if val := readLine(reader); val != "" {
		cfg.DefaultProvider = val
	}

	fmt.Println("\n--- Claude Settings ---")
	fmt.Printf("Claude model [%s] (e.g. claude-3-5-sonnet-20240620): ", cfg.Providers.Claude.Model)
	if val := readLine(reader); val != "" {
		cfg.Providers.Claude.Model = val
	}

	fmt.Println("\n--- GitHub Models Settings ---")
	fmt.Printf("GitHub token [%s]: ", maskToken(cfg.Providers.GitHub.Token))
	if val := readLine(reader); val != "" {
		cfg.Providers.GitHub.Token = val
	}
	fmt.Printf("GitHub model [%s] (e.g. gpt-4o): ", cfg.Providers.GitHub.Model)
	if val := readLine(reader); val != "" {
		cfg.Providers.GitHub.Model = val
	}

	fmt.Println("\n--- Ollama Settings ---")
	fmt.Printf("Ollama host [%s] (e.g. http://localhost:11434): ", cfg.Providers.Ollama.Host)
	if val := readLine(reader); val != "" {
		cfg.Providers.Ollama.Host = val
	}
	fmt.Printf("Ollama model [%s] (e.g. llama3.2): ", cfg.Providers.Ollama.Model)
	if val := readLine(reader); val != "" {
		cfg.Providers.Ollama.Model = val
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Printf("\nConfig saved to %s\n", config.Path())
	return nil
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func maskToken(token string) string {
	if token == "" {
		return "(not set)"
	}
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
