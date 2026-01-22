package promptbox

import (
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

		case tea.KeyEsc:
			// This removes focus from the text editing
			// widget: useful for when we want to scroll
			// up to view messages logged to the console.
			if m.textarea.Focused() {
				m.textarea.Blur()
			}

		case tea.KeyRunes:
			// Whenever the user starts typing _letters_,
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

	view.WriteString("\n" + m.textarea.View())
	view.WriteString("\n\n(ctrl+c to quit, ctrl+\\ to submit prompt)\n\n")

	return view.String()
}

func GetPrompt() (string, error) {
	txt := textarea.New()
	txt.Placeholder = "Start typing your prompt to focus here, or Esc to remove focus"
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
