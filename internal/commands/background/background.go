package background

import (
	"t-kt/internal/commands"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type BgTaskResult struct {
	Result commands.CmdResult
	TaskId int
}

type BgTask func() BgTaskResult

type BgTaskManager []BgTask

func NewBgTaskManager(cmds ...func() commands.CmdResult) BgTaskManager {
	return bgWrap(cmds...)
}

func (bg BgTaskManager) GetTask(taskId int) BgTask {
	return bg[taskId]
}

func (bg BgTaskManager) GetTasks() []BgTask {
	return bg
}

func (bg BgTaskManager) GetWrapTask(taskId int) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second / 2)
		return bg[taskId]()
	}
}

func (bg BgTaskManager) GetWrapTasks() []tea.Cmd {
	tasks := make([]tea.Cmd, 0, len(bg))

	for _, cmd := range bg {
		wrapped := func() tea.Msg {
			time.Sleep(time.Second / 2)
			return cmd()
		}

		tasks = append(tasks, wrapped)
	}
	return tasks
}

func bgWrap(cmds ...func() commands.CmdResult) []BgTask {
	bgTasks := make(BgTaskManager, 0, len(cmds))

	for idx, cmd := range cmds {
		wrapped := func() BgTaskResult {
			res := cmd()
			return BgTaskResult{
				Result: res,
				TaskId: idx,
			}
		}
		bgTasks = append(bgTasks, wrapped)
	}

	return bgTasks
}
