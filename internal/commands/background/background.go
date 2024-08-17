package background

import "t-kt/internal/commands"

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
