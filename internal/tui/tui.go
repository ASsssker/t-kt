package tui

import (
	"t-kt/internal/commands"

	tea "github.com/charmbracelet/bubbletea"
)

func NewTUI() tea.Model {
	options := []tea.Model{
		newButton("Перезапустить сервер", wrap(commands.RestartServer)),
		newButton("Запустить сервер", wrap(commands.StartServer)),
		newButton("Остановить сервер", wrap(commands.StopServer)),
		newButton("Собрать саппорт", wrap(commands.ExtractLogs)),
		newButton("Включить дебаг для сервера", wrap(commands.SwitchToDebug)),
		newButton("Закрыть клиент", wrap(commands.KillUI)),
	}
	section1 := newSection(options)
	screen := newScreen([]tea.Model{section1}, []string{"AN"})

	return screen
}

func wrap(cmd func() commands.CmdResult) func() tea.Msg {
	return func() tea.Msg {
		return cmd()
	}
}
