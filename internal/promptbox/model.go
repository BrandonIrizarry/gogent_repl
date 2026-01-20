package promptbox

import (
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	textarea   textarea.Model
	err        error
	promptText string
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			slog.Debug("Model:", slog.Int("prompt_text_length", len(m.promptText)))
			return m, tea.Quit

		case tea.KeyCtrlBackslash:
			m.promptText = m.textarea.Value()
			return m, tea.Quit

		default:
			// Any keystroke other than quitting will
			// bring focus back to the text editing
			// widget.
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var view strings.Builder

	view.WriteString(m.textarea.View())
	view.WriteString("\n\n(ctrl+c to quit, ctrl+\\ to submit prompt)\n\n")

	return view.String()
}
