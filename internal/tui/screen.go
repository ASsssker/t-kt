package tui

import (
	"fmt"
	"strings"
	"t-kt/internal/commands"
	"t-kt/internal/commands/background"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type Screen struct {
	sections       []tea.Model
	sectionsName   []string
	currentSection int
	isLoaded       bool
	warnMsg        []string
	bgManager      background.BgTaskManager
	loader         loader
	width          int
	height         int
	selectedStyle  lipgloss.Style
	helpStyle      lipgloss.Style
}

func newScreen(sections []tea.Model, sectionsName []string, bgTasks ...func() commands.CmdResult) Screen {
	bgManager := background.NewBgTaskManager(bgTasks...)
	l := newLoader([]string{`\`, "|", "/", "―"})

	return Screen{
		sections:      sections,
		sectionsName:  sectionsName,
		bgManager:     bgManager,
		loader:        l,
		selectedStyle: newSelectedStyle(),
		helpStyle:     newHelpStyle(),
	}
}

func (screen Screen) Init() tea.Cmd {
	cmds := append(screen.bgManager.GetWrapTasks(), tea.ClearScreen, screen.loader.Tick)
	return tea.Batch(cmds...)
}

func (screen Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	// Обработка результатов выполнения команд
	case commands.CmdResult:
		screen.isLoaded = false
		return screen.handleCmdResult(msg)
	// Обработка фоновых задач
	case background.BgTaskResult:
		return screen.handleBg(msg)
	// Обработка анимации загрузки
	case loader:
		screen.loader = msg
		return screen, screen.loader.Tick
	// Обработка событий изменения размера консоли
	case tea.WindowSizeMsg:
		screen.width = msg.Width
		screen.height = msg.Height
		return screen, nil
	// Обработка нажатий кнопок не связанных с изменением интерфейса
	case tea.KeyMsg:
		switch msg.String() {
		// Очистка вывода
		case "c":
			screen.warnMsg = []string{}
			return screen, nil
		}
	}

	// Блокируем интерфейс во время выпололнения команды
	if screen.isLoaded {
		return screen, nil
	}

	var cmd tea.Cmd
	var model tea.Model

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Обработка нажатий кнопок связанных с изменением интерфейса
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
		if idx == screen.currentSection {
			sectionNavBar += screen.selectedStyle.Render(fmt.Sprintf("%s\t", screen.sectionsName[idx]))
			continue
		}

		sectionNavBar += fmt.Sprintf("%s\t", screen.sectionsName[idx])
	}

	sectionNavBar += "\n\n"
	sectionNavBar += screen.sections[screen.currentSection].View()
	sectionNavBar += screen.viewWarnMsg()

	sectionNavBar = screen.viewHelp(sectionNavBar)

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
	if str := msg.GetMsg(); str != "" {
		screen.warnMsg = append(screen.warnMsg, str)
	}
	return screen, nil
}

func (screen Screen) handleBg(msg background.BgTaskResult) (tea.Model, tea.Cmd) {
	if str := msg.Result.GetMsg(); str != "" {
		screen.warnMsg = append(screen.warnMsg, str)
	}

	return screen, screen.bgManager.GetWrapTask(msg.TaskId)
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
		if len(msg) >= screen.width-5 {
			warnMsgs += msg[:screen.width-5] + "...\n"
		} else {
			warnMsgs += fmt.Sprintf("%s\n", msg)
		}
	}

	// Если сообщений меньше msgCount добавляем пустые строки
	for i := warnMsgLen; i < msgCount; i++ {
		warnMsgs += "\n"
	}

	return warnMsgs
}

func (screen Screen) viewHelp(currentView string) string {
	rowCount := strings.Count(currentView, "\n")
	for i := rowCount; i <= screen.height-5; i++ {
		currentView += "\n"
	}
	if screen.isLoaded {
		currentView += screen.loader.View() + screen.helpStyle.Render("\nc очистка вывода\tq выход \n")
	} else {
		currentView += screen.helpStyle.Render("\nc очистка вывода\tq выход \n")

	}
	return currentView
}

func newSelectedStyle() lipgloss.Style {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#7FFF00"))
	return style
}

func newHelpStyle() lipgloss.Style {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#787878")).Align(lipgloss.Right).Width(80)
	return style
}
