package tui

import (
	"fmt"
	"t-kt/internal/commands"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
)

type Screen struct {
	sections       []tea.Model
	sectionsName   []string
	currentSection int
	isLoaded       bool
	warnMsg        []string
	width          int
	height         int
}

func (screen Screen) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen)
}

func (screen Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case commands.CmdResult:
		screen.isLoaded = false
		return screen.handleCmdResult(msg)
	}

	// Блокируем интерфейс во время выпололнения команды
	if screen.isLoaded {
		return screen, nil
	}

	var cmd tea.Cmd
	var model tea.Model

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		screen.width = msg.Width
		screen.height = msg.Height
		return screen, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "a", "d", "left", "right":
			model, cmd = screen.changeSection(msg)
			return model, cmd
		case "q", "ctrl+c":
			return screen, tea.Quit
		case "w", "s", "up", "down":
			return screen.updateSection(msg)
		case "enter", " ":
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
	sectionNavBar += screen.viewWarnMsg()

	sectionNavBar += "\nНажмите q для выхода.\n"

	return wordwrap.String(sectionNavBar, screen.width)
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

func (screen Screen) handleCmdResult(msg commands.CmdResult) (tea.Model, tea.Cmd) {
	if msg != nil {
		screen.warnMsg = append(screen.warnMsg, msg.Error())
	}

	return screen, nil
}

func (screen Screen) viewWarnMsg() string {
	var msgSlice []string
	warnMsgs := "\n"
	msgCount := 5

	warnMsgLen := len(screen.warnMsg)

	if warnMsgLen < msgCount {
		msgSlice = screen.warnMsg
	} else {
		msgSlice = screen.warnMsg[warnMsgLen-msgCount:]
	}

	for _, msg := range msgSlice {
		warnMsgs += fmt.Sprintf("%s\n", msg)
	}

	for i := warnMsgLen; i < msgCount; i++ {
		warnMsgs += "\n"
	}

	return warnMsgs
}
