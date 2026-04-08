package cmd

import (
	"fmt"

	"github.com/WMcKibbin/howdoi/internal/shell"
	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Generate shell aliases for howdoi",
	Long: `Generate shell-specific aliases for howdoi.

Add the output to your shell config file:
  bash: eval "$(howdoi alias)"
  zsh:  eval "$(howdoi alias)"
  fish: howdoi alias | source
  pwsh: howdoi alias | Invoke-Expression`,
	RunE: runAlias,
}

var aliasShell string

func init() {
	aliasCmd.Flags().StringVar(&aliasShell, "shell", "", "Shell type (bash, zsh, fish, powershell). Auto-detected if not set.")
	rootCmd.AddCommand(aliasCmd)
}

func runAlias(cmd *cobra.Command, args []string) error {
	sh := aliasShell
	if sh == "" {
		info := shell.Detect()
		sh = info.Shell
	}
	fmt.Println(shell.GenerateAlias(sh))
	return nil
}
