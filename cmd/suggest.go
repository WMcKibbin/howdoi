package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/WMcKibbin/howdoi/internal/prompt"
	"github.com/WMcKibbin/howdoi/internal/shell"
	"github.com/WMcKibbin/howdoi/internal/ui"
	"github.com/spf13/cobra"
)

var suggestCmd = &cobra.Command{
	Use:   "suggest [query]",
	Short: "Get a command suggestion for a task",
	Long:  "Describe what you want to do and get a shell command suggestion from an AI provider.",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSuggest,
}

func init() {
	rootCmd.AddCommand(suggestCmd)
}

func runSuggest(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")
	p, err := getProvider()
	if err != nil {
		return err
	}

	info := shell.Detect()
	systemPrompt := prompt.Suggest(info.OS, info.Shell)

	ctx := context.Background()
	suggestion, err := p.Chat(ctx, systemPrompt, query)
	if err != nil {
		return fmt.Errorf("provider %s: %w", p.Name(), err)
	}

	// Clean up any markdown code fences the model might have included
	suggestion = cleanCommand(suggestion)

	return interactiveLoop(ctx, p, info, query, suggestion)
}

func interactiveLoop(ctx context.Context, p interface {
	Chat(ctx context.Context, systemPrompt string, userMessage string) (string, error)
	Name() string
}, info shell.Info, originalQuery, command string) error {
	for {
		action, err := ui.ShowMenu(command)
		if err != nil {
			return err
		}

		switch action {
		case ui.ActionExecute:
			confirmed, err := ui.Confirm("Are you sure you want to execute the suggested command?")
			if err != nil {
				return err
			}
			if confirmed {
				fmt.Println()
				if err := shell.Execute(command); err != nil {
					return err
				}
				_ = info.AppendToHistory(command)
			}
			return nil

		case ui.ActionCopy:
			if err := shell.CopyToClipboard(command); err != nil {
				return fmt.Errorf("copying to clipboard: %w", err)
			}
			fmt.Println("  Copied to clipboard!")
			return nil

		case ui.ActionRevise:
			revision, err := ui.PromptInput("How would you like to revise this command?")
			if err != nil {
				return err
			}
			if revision == "" {
				continue
			}
			revisedQuery := fmt.Sprintf("Original request: %s\nOriginal command: %s\nRevision: %s", originalQuery, command, revision)
			systemPrompt := prompt.Suggest(info.OS, info.Shell)
			newCmd, err := p.Chat(ctx, systemPrompt, revisedQuery)
			if err != nil {
				return fmt.Errorf("provider %s: %w", p.Name(), err)
			}
			command = cleanCommand(newCmd)

		case ui.ActionExplain:
			explainPrompt := prompt.Explain()
			explanation, err := p.Chat(ctx, explainPrompt, command)
			if err != nil {
				return fmt.Errorf("provider %s: %w", p.Name(), err)
			}
			fmt.Printf("\n%s\n", explanation)
			return nil

		case ui.ActionCancel:
			return nil
		}
	}
}

func cleanCommand(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		lines := strings.Split(s, "\n")
		if len(lines) >= 3 {
			lines = lines[1 : len(lines)-1]
			s = strings.Join(lines, "\n")
		}
	}
	return strings.TrimSpace(s)
}
