package tui

import (
	"fmt"
	"t-kt/internal/commands"
	"t-kt/internal/commands/background"
	"t-kt/internal/configs"

	tea "github.com/charmbracelet/bubbletea"
)

func NewTUI(conf configs.Config) tea.Model {
	options1 := []tea.Model{
		newButton("Очистить логи", wrap(commands.ClearLogs)),
		newButton("Перезапустить сервер", wrap(commands.RestartServer)),
		newButton("Запустить сервер", wrap(commands.StartServer)),
		newButton("Остановить сервер", wrap(commands.StopServer)),
		newButton("Собрать саппорт", wrap(commands.ExtractLogs)),
		newButton("Включить дебаг для сервера", wrap(commands.SwitchToDebug)),
		newButton("Закрыть клиент", wrap(commands.KillUI)),
		newButton("Чистка обсолетов", wrap(commands.DisableObselete)),
	}
	var test bool

	switcher, _ := commands.NewArchiveSwitcher(&test, conf.IPC)
	fmt.Println("dsd")
	options2 := []tea.Model{ newCheckbox("Запись отрезками", switcher, nil)}
	section1 := newSection(options1)
	section2 := newSection(options2)
	screen := newScreen([]tea.Model{section1, section2}, []string{"AN", "RS"}, background.CheckDump)

	return screen
}

func wrap(cmd func() commands.CmdResult) func() tea.Msg {
	return func() tea.Msg {
		return cmd()
	}
}
