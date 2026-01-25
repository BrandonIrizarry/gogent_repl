package radioselect

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursorPos int
	choice    string
	choices   []string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEnter:
			m.choice = m.choices[m.cursorPos]
			return m, tea.Quit

		case tea.KeyDown:
			m.cursorPos++

			// Wrap back to the top if necessary.
			if m.cursorPos >= len(m.choices) {
				m.cursorPos = 0
			}

		case tea.KeyUp:
			m.cursorPos--

			// Wrap back to the bottom if necessary.
			if m.cursorPos < 0 {
				m.cursorPos = len(m.choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString("\nRecent projects\n\n")

	for i := range len(m.choices) {
		if m.cursorPos == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(m.choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(Press Enter to submit, ctrl+c to quit)\n\n")

	return s.String()
}

func SelectWorkingDir(choices []string) (string, error) {
	p := tea.NewProgram(model{
		choices: choices,
	})

	m, err := p.Run()
	if err != nil {
		return "", err
	}

	rsel := m.(model)
	return rsel.choice, nil
}
