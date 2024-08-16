package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type button struct {
	title string
	f     func() tea.Msg
	style lipgloss.Style
}

func newButton(title string, f func() tea.Msg) button {
	return button{
		title: title,
		f:     f,
		style: newButtonStyle(),
	}
}

func (b button) Init() tea.Cmd {
	return nil
}

func (b button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			return b.action()
		}
	}

	return b, nil
}

func (b button) View() string {
	return b.style.Render(b.title)
}

func (b button) action() (tea.Model, tea.Cmd) {
	// ...
	return b, b.f
}

func newButtonStyle() lipgloss.Style {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFE4C4"))
	return style
}
