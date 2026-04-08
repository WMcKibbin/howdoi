package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/WMcKibbin/howdoi/internal/prompt"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain [command]",
	Short: "Explain what a shell command does",
	Long:  "Provide a shell command and get a detailed explanation of what it does, including flags and arguments.",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runExplain,
}

func init() {
	rootCmd.AddCommand(explainCmd)
}

func runExplain(cmd *cobra.Command, args []string) error {
	command := strings.Join(args, " ")
	p, err := getProvider()
	if err != nil {
		return err
	}

	systemPrompt := prompt.Explain()
	ctx := context.Background()
	explanation, err := p.Chat(ctx, systemPrompt, command)
	if err != nil {
		return fmt.Errorf("provider %s: %w", p.Name(), err)
	}

	fmt.Println(explanation)
	return nil
}
