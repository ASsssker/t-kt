package commands

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

func ClearLogs() CmdResult {
	err := clearDir(ServerLogAbsPath)
	if err != nil {
		return NewCmdResult("", err)
	}

	clientLogDir, err := GetClientLogDir()
	if err != nil {
		return NewCmdResult("", err)
	}

	err = clearDir(clientLogDir)
	if err != nil {
		return NewCmdResult("", err)
	}

	return NewCmdResult("", nil)
}

func ExtractLogs() CmdResult {
	cmd := exec.Command(extractLogAppPath, supportDstPath)
	if err := cmd.Run(); err != nil {
		return NewCmdResult("", err)
	}

	return NewCmdResult("", nil)
}

func KillUI() CmdResult {
	cmd := exec.Command("pidof", clientProcessName)
	out, err := cmd.Output()
	if err != nil {
		return NewCmdResult("", err)
	}

	pid := bytes.TrimSpace(out)
	cmd = exec.Command("kill", string(pid))
	if err = cmd.Run(); err != nil {
		return NewCmdResult("", err)
	}

	return NewCmdResult("", nil)
}

func RestartServer() CmdResult {
	cmd := exec.Command("service", strings.Split(restartServer, " ")...)
	if err := cmd.Run(); err != nil {
		return NewCmdResult("", err)
	}

	return NewCmdResult("", nil)
}

func StartServer() CmdResult {
	cmd := exec.Command("service", strings.Split(startServer, " ")...)
	if err := cmd.Run(); err != nil {
		return NewCmdResult("", err)
	}

	return NewCmdResult("", nil)
}

func StopServer() CmdResult {
	cmd := exec.Command("service", strings.Split(stopServer, " ")...)
	if err := cmd.Run(); err != nil {
		return NewCmdResult("", err)
	}

	return NewCmdResult("", nil)
}

func SwitchToDebug() CmdResult {
	file, err := os.OpenFile(serverConfFilePath, os.O_RDWR, 0)
	if err != nil {
		return NewCmdResult("", err)
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return NewCmdResult("", err)
	}

	buf = bytes.Replace(buf, []byte("INFO"), []byte("DEBUG"), 1)

	if err = file.Truncate(0); err != nil {
		return CmdResult{err: err}
	}
	if _, err = file.Seek(0, 0); err != nil {
		return CmdResult{err: err}
	}

	if _, err = file.Write(buf); err != nil {
		return CmdResult{err: err}
	}

	return NewCmdResult("", nil)
}
