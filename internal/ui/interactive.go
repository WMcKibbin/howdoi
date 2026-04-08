package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Action represents what the user chose from the interactive menu.
type Action int

const (
	ActionExecute Action = iota
	ActionRevise
	ActionExplain
	ActionCopy
	ActionCancel
)

var choices = []string{
	"Execute command",
	"Revise command",
	"Explain command",
	"Copy to clipboard",
	"Cancel",
}

var actions = []Action{
	ActionExecute,
	ActionRevise,
	ActionExplain,
	ActionCopy,
	ActionCancel,
}

type model struct {
	cursor  int
	command string
	chosen  Action
	done    bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(choices)-1 {
				m.cursor++
			}
		case "enter":
			m.chosen = actions[m.cursor]
			m.done = true
			return m, tea.Quit
		case "q", "ctrl+c", "esc":
			m.chosen = ActionCancel
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	commandStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).PaddingLeft(4)
	promptStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
)

func (m model) View() string {
	if m.done {
		return ""
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("  Suggestion:"))
	b.WriteString("\n\n")
	b.WriteString(commandStyle.Render(m.command))
	b.WriteString("\n\n")
	b.WriteString(promptStyle.Render("  ? Select an option"))
	b.WriteString("\n")

	for i, choice := range choices {
		if m.cursor == i {
			b.WriteString(cursorStyle.Render("  > " + choice))
		} else {
			b.WriteString("    " + choice)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// ShowMenu displays the interactive post-suggestion menu and returns the chosen action.
func ShowMenu(command string) (Action, error) {
	m := model{command: command}
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return ActionCancel, err
	}
	return result.(model).chosen, nil
}

// Confirm shows a yes/no confirmation prompt and returns the result.
func Confirm(prompt string) (bool, error) {
	m := confirmModel{prompt: prompt}
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return false, err
	}
	return result.(confirmModel).confirmed, nil
}

type confirmModel struct {
	prompt    string
	cursor    int
	confirmed bool
	done      bool
}

func (m confirmModel) Init() tea.Cmd { return nil }

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "left", "h":
			m.cursor = 0
		case "down", "j", "right", "l":
			m.cursor = 1
		case "enter":
			m.confirmed = m.cursor == 0
			m.done = true
			return m, tea.Quit
		case "q", "ctrl+c", "esc":
			m.confirmed = false
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	if m.done {
		return ""
	}
	opts := []string{"Yes", "No"}
	var b strings.Builder
	b.WriteString(promptStyle.Render(fmt.Sprintf("  ? %s", m.prompt)))
	b.WriteString("\n")
	for i, opt := range opts {
		if m.cursor == i {
			b.WriteString(cursorStyle.Render("  > " + opt))
		} else {
			b.WriteString("    " + opt)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// PromptInput shows a simple text input prompt and returns what the user typed.
func PromptInput(promptText string) (string, error) {
	m := inputModel{prompt: promptText}
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	return result.(inputModel).value, nil
}

type inputModel struct {
	prompt string
	value  string
	done   bool
}

func (m inputModel) Init() tea.Cmd { return nil }

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.done = true
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.value = ""
			m.done = true
			return m, tea.Quit
		case "backspace":
			if len(m.value) > 0 {
				m.value = m.value[:len(m.value)-1]
			}
		default:
			if len(msg.String()) == 1 || msg.String() == " " {
				m.value += msg.String()
			}
		}
	}
	return m, nil
}

func (m inputModel) View() string {
	if m.done {
		return ""
	}
	return fmt.Sprintf("%s\n  > %s\n", promptStyle.Render(fmt.Sprintf("  ? %s", m.prompt)), m.value)
}
