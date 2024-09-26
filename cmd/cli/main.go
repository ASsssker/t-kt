package main

import (
	"flag"
	"log"
	"t-kt/internal/configs"
	"t-kt/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var IPCConf configs.IPCConf
	flag.StringVar(&IPCConf.DNS, "ipc", "", "Данные для входа в камеру в формате username:password@host")
	flag.IntVar(&IPCConf.ArchiveSwitchTime, "archive_switch_time", 0, "Время которое проходит между вкл/выкл архива в секундах. По умолчанию рандомное значение от 0 до 30")
	flag.Parse()
	conf := configs.Config{
		IPC: IPCConf,
	}

	t := tea.NewProgram(tui.NewTUI(conf))
	if _, err := t.Run(); err != nil {
		log.Fatal(err)
	}
}
