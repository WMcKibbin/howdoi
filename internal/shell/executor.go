package shell

import (
	"fmt"
	"os"
	"os/exec"
)

// Execute runs a command in the user's shell.
func Execute(command string) error {
	info := Detect()
	var cmd *exec.Cmd

	switch info.Shell {
	case "powershell":
		cmd = exec.Command("powershell", "-Command", command)
	case "cmd":
		cmd = exec.Command("cmd", "/C", command)
	default:
		shellBin := os.Getenv("SHELL")
		if shellBin == "" {
			shellBin = "/bin/sh"
		}
		cmd = exec.Command(shellBin, "-c", command)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	return nil
}
