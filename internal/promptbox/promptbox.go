package promptbox

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func GetPrompt() (string, error) {
	txt := textarea.New()
	txt.Placeholder = "Write your prompt here..."
	txt.ShowLineNumbers = false

	// Don't set a background highlighting color, since at least
	// in Vterm under an Emacs light theme this will still paint
	// the background black (lipgloss thinks the terminal is using
	// a dark theme somehow.)
	txt.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(lipgloss.Color(""))
	p := tea.NewProgram(model{textarea: txt})

	m, err := p.Run()
	if err != nil {
		return "", err
	}

	pbox := m.(model)
	return pbox.promptText, nil
}
