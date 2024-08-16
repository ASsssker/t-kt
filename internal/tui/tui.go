package tui

import (
	"t-kt/internal/commands"

	tea "github.com/charmbracelet/bubbletea"
)

func NewTUI() tea.Model {
	options := []tea.Model{
		button{
			title: "Перезапустить сервер",
			f:     wrap(commands.RestartServer),
		},
		button{
			title: "Запустить сервер",
			f:     wrap(commands.StartServer),
		},
		button{
			title: "Остановить сервер",
			f:     wrap(commands.StopServer),
		},
		button{
			title: "Собрать саппорт",
			f:     wrap(commands.ExtractLogs),
		},
		button{
			title: "Включить дебаг для сервера",
			f:     wrap(commands.SwitchToDebug),
		},
		button{
			title: "Закрыть клиент",
			f:     wrap(commands.KillUI),
		},
	}

	section1 := section{options: options}
	screen := Screen{
		sections:     []tea.Model{section1},
		sectionsName: []string{"AN"},
	}

	return screen
}

func wrap(cmd func() commands.CmdResult) func() tea.Msg {
	return func() tea.Msg {
		return cmd()
	}
}
