package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loader struct {
	icon  []string
	t     int
	style lipgloss.Style
}

func newLoader(icon []string) loader {
	return loader{
		icon:  icon,
		style: newLoaderStyle(),
	}
}

func (l loader) Tick() tea.Msg {
	time.Sleep(time.Second / 10)
	if l.t >= len(l.icon) {
		l.t = 0
	} else {
		l.t++
	}
	return l
}

func (l loader) View() string {
	idx := l.t % len(l.icon)
	return l.style.Render(l.icon[idx])
}

func newLoaderStyle() lipgloss.Style {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFE0"))
	return style
}
