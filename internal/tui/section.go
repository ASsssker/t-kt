package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type section struct {
	options     []tea.Model
	cursor      int
	cursorStyle lipgloss.Style
}

func newSection(options []tea.Model) section {
	return section{
		options:     options,
		cursorStyle: newCursorStyle(),
	}
}

func (s section) Init() tea.Cmd {
	return nil
}

func (s section) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "w", "s", "up", "down":
			return s.moveCursor(msg)
		case "enter", " ":
			return s.updateOption(msg)
		}
	}

	return s, nil
}

func (s section) View() string {
	var optionList string

	for idx := range s.options {
		cursor := " "
		if idx == s.cursor {
			cursor = s.cursorStyle.Render(">")
		}

		optionList += fmt.Sprintf("%s %s\n", cursor, s.options[idx].View())
	}

	return optionList
}

func (s section) moveCursor(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "w", "up":
		if s.cursor > 0 {
			s.cursor--
		}

		return s, nil

	case "s", "down":
		if s.cursor < len(s.options)-1 {
			s.cursor++
		}

		return s, nil
	default:
		return s, nil
	}
}

func (s section) updateOption(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	s.options[s.cursor], cmd = s.options[s.cursor].Update(msg)
	return s, cmd
}

func newCursorStyle() lipgloss.Style {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#7FFF00"))
	return style
}
