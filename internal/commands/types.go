package commands

type CmdResult struct {
	info string
	err error
}

func(c CmdResult) GetMsg() string {
	if c.err != nil {
		return c.err.Error()
	}

	return c.info
}


