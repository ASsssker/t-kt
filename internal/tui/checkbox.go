package tui

import (
	"context"
	"sync"
	"t-kt/internal/commands"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type checkBoxMsg string

type checkbox struct {
	title        string
	status       bool
	enableF      func() tea.Msg
	disableF     func() tea.Msg
	defaultStyle lipgloss.Style
	checkedStyle lipgloss.Style
}

func newCheckbox(title string, enableF func(context.Context) commands.CmdResult, disableF func() commands.CmdResult) checkbox {
	enable, disable := func() (func(), func() tea.Msg) {
		var mu sync.Mutex
		var ctx context.Context
		var cancel context.CancelFunc

		e := func() {
			mu.Lock()
			defer mu.Unlock()
			if disableF == nil {
				ctx, cancel = context.WithCancel(context.Background())
				go enableF(ctx)
			} else {
				enableF(nil)
			}
		}

		c := func() tea.Msg {
			mu.Lock()
			defer mu.Unlock()
			if cancel != nil {
				cancel()
				return checkBoxMsg("Задача остановлена")
			}
			disableF()
			return checkBoxMsg("Задача остановлена")
		}

		return e, c
	}()

	return checkbox{
		title:        title,
		enableF:      checkboxWrap(enable),
		disableF:     disable,
		defaultStyle: newButtonStyle(),
		checkedStyle: newSelectedStyle(),
	}
}

func (c checkbox) Init() tea.Cmd {
	return nil
}

func (c checkbox) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			if c.status {
				c.status = false
				return c, c.disableF
			} else {
				c.status = true
				return c, c.enableF
			}
		}
	}

	return c, nil
}

func (c checkbox) View() string {
	if c.status {
		return c.checkedStyle.Render("[X]" + c.title)
	}

	return c.defaultStyle.Render("[ ]" + c.title)
}

func checkboxWrap(f func()) func() tea.Msg {
	return func() tea.Msg {
		go f()
		return checkBoxMsg("Задача запущена")
	}
}
