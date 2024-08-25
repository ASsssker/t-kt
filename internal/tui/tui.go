package tui

import (
	"fmt"
	"t-kt/internal/commands"
	"t-kt/internal/commands/background"

	tea "github.com/charmbracelet/bubbletea"
)

func NewTUI() tea.Model {
	options1 := []tea.Model{
		newButton("Очистить логи", wrap(commands.ClearLogs)),
		newButton("Перезапустить сервер", wrap(commands.RestartServer)),
		newButton("Запустить сервер", wrap(commands.StartServer)),
		newButton("Остановить сервер", wrap(commands.StopServer)),
		newButton("Собрать саппорт", wrap(commands.ExtractLogs)),
		newButton("Включить дебаг для сервера", wrap(commands.SwitchToDebug)),
		newButton("Закрыть клиент", wrap(commands.KillUI)),
	}
	options2 := []tea.Model{
		newCheckbox("тест1", func() tea.Msg { fmt.Print("Hi"); return "" }, func() tea.Msg { return 0 }, true),
		newCheckbox("тест2", func() tea.Msg { fmt.Print("Do"); return "" }, func() tea.Msg { fmt.Print("Do cancel"); return 0 }, false),
	}
	section1 := newSection(options1)
	section2 := newSection(options2)
	screen := newScreen([]tea.Model{section1, section2}, []string{"AN", "Testing"}, background.CheckDump)

	return screen
}

func wrap(cmd func() commands.CmdResult) func() tea.Msg {
	return func() tea.Msg {
		return cmd()
	}
}
