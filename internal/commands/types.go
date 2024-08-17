package commands

type CmdResult struct {
	info string
	err  error
}

func NewCmdResult(info string, err error) CmdResult {
	return CmdResult{
		info: info,
		err:  err,
	}
}

func (c CmdResult) GetMsg() string {
	if c.err != nil {
		return c.err.Error()
	}

	return c.info
}
