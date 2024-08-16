package main

import (
	"log"
	"t-kt/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	t := tea.NewProgram(tui.NewTUI())
	if _, err := t.Run(); err != nil {
		log.Fatal(err)
	}
}
