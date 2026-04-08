package shell

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Info holds detected shell environment information.
type Info struct {
	OS    string // "macos", "linux", "windows"
	Shell string // "bash", "zsh", "fish", "powershell", "unknown"
}

// Detect returns the current OS and shell.
func Detect() Info {
	info := Info{}

	switch runtime.GOOS {
	case "darwin":
		info.OS = "macos"
	case "linux":
		info.OS = "linux"
	case "windows":
		info.OS = "windows"
	default:
		info.OS = runtime.GOOS
	}

	if runtime.GOOS == "windows" {
		if os.Getenv("PSModulePath") != "" {
			info.Shell = "powershell"
		} else {
			info.Shell = "cmd"
		}
	} else {
		shellPath := os.Getenv("SHELL")
		info.Shell = filepath.Base(shellPath)
		if info.Shell == "" || info.Shell == "." {
			info.Shell = "unknown"
		}
	}

	return info
}

// HistoryFile returns the path to the shell history file, if known.
func (i Info) HistoryFile() string {
	home, _ := os.UserHomeDir()
	if home == "" {
		return ""
	}
	switch i.Shell {
	case "bash":
		return filepath.Join(home, ".bash_history")
	case "zsh":
		if f := os.Getenv("HISTFILE"); f != "" {
			return f
		}
		return filepath.Join(home, ".zsh_history")
	case "fish":
		return filepath.Join(home, ".local", "share", "fish", "fish_history")
	default:
		return ""
	}
}

// AppendToHistory appends a command to the shell history file.
func (i Info) AppendToHistory(command string) error {
	histFile := i.HistoryFile()
	if histFile == "" {
		return nil
	}

	f, err := os.OpenFile(histFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	entry := command
	if i.Shell == "zsh" {
		entry = ": 0:0;" + command
	}
	if !strings.HasSuffix(entry, "\n") {
		entry += "\n"
	}
	_, err = f.WriteString(entry)
	return err
}
