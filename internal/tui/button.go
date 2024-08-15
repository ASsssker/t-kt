package tui

import tea "github.com/charmbracelet/bubbletea"

type button struct {
	title string
	f     func() tea.Msg
}

func (b button) Init() tea.Cmd {
	return nil
}

func (b button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "space":
			return b.action()
		}
	}

	return b, nil
}

func (b button) View() string {
	return b.title
}

func (b button) action() (tea.Model, tea.Cmd) {
	// ...
	return b, b.f
}
