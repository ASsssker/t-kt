package tui

import (
	"log"
	"t-kt/internal/commands"
	"t-kt/internal/commands/background"
	"t-kt/internal/configs"

	tea "github.com/charmbracelet/bubbletea"
)

func NewTUI(conf configs.Config) tea.Model {
	sections := []tea.Model{}
	section_name := []string{"AN"}

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
	sections = append(sections, newSection(options1))

	if len(conf.IPC.DNS) != 0 {
		switcher, err := commands.NewArchiveSwitcher(conf.IPC)
		if err != nil {
			log.Fatal(err)
		}
		options2 := []tea.Model{newCheckbox("Запись отрезками", switcher, nil)}
		sections = append(sections, newSection(options2))
		section_name = append(section_name, "RS")

	}

	screen := newScreen(sections, section_name, background.CheckDump)

	return screen
}

func wrap(cmd func() commands.CmdResult) func() tea.Msg {
	return func() tea.Msg {
		return cmd()
	}
}
