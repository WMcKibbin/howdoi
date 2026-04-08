package shell

import "fmt"

// GenerateAlias returns shell-specific alias configuration.
func GenerateAlias(shellName string) string {
	switch shellName {
	case "bash", "zsh":
		return "alias hdi='howdoi suggest'\nalias hde='howdoi explain'"
	case "fish":
		return "alias hdi 'howdoi suggest'\nalias hde 'howdoi explain'"
	case "powershell":
		return "function hdi { howdoi suggest @args }\nfunction hde { howdoi explain @args }"
	default:
		return fmt.Sprintf("# Unsupported shell: %s\n# Add manually:\n#   alias hdi='howdoi suggest'\n#   alias hde='howdoi explain'", shellName)
	}
}
