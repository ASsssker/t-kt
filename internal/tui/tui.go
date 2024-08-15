package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Screen struct {
	sections       []tea.Model
	sectionsName   []string
	currentSection int
	isLoaded       bool
	warnMsg        []string
}

func (screen Screen) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen)
}

func (screen Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Блокируем интерфейс во время выпололнения команды
	if screen.isLoaded {
		return screen, nil
	}

	var cmd tea.Cmd
	var model tea.Model

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a", "d", "left", "right":
			model, cmd = screen.changeSection(msg)
			return model, cmd
		case "q", "ctrl+c":
			return screen, tea.Quit
		case "w", "s", "up", "down":
			return screen.updateSection(msg)
		case "enter", "space":
			screen.isLoaded = true
			return screen.updateSection(msg)
		}
	default:
		return screen, nil
	}
	return screen, nil
}

func (screen Screen) View() string {
	var sectionNavBar string

	for idx := range screen.sections {
		selected := " "
		if idx == screen.currentSection {
			selected = "+"
		}

		sectionNavBar += fmt.Sprintf("[%s] %s\t", selected, screen.sectionsName[idx])
	}

	sectionNavBar += "\n"
	sectionNavBar += screen.sections[screen.currentSection].View()
	sectionNavBar += "\nНажмите q для выхода.\n"

	return sectionNavBar
}

func (screen Screen) changeSection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "a", "left":
		if screen.currentSection > 0 {
			screen.currentSection--
		}

		return screen, nil

	case "d", "right":
		if screen.currentSection < len(screen.sections)-1 {
			screen.currentSection++
		}

		return screen, nil
	default:
		return screen, nil
	}
}

func (screen Screen) updateSection(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	screen.sections[screen.currentSection], cmd = screen.sections[screen.currentSection].Update(msg)
	return screen, cmd
}
