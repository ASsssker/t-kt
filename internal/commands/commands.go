package commands

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

func ClearLogs() CmdResult {
	err := clearDir(serverLogAbsPath)
	if err != nil {
		return CmdResult{err: err}
	}

	clientLogDir, err := getClientLogDir()
	if err != nil {
		return CmdResult{err: err}
	}

	err = clearDir(clientLogDir)
	if err != nil {
		return CmdResult{err: err}
	}

	return CmdResult{}
}

func ExtractLogs() CmdResult {
	cmd := exec.Command(extractLogAppPath, supportDstPath)
	if err := cmd.Run(); err != nil {
		return CmdResult{err: err}
	}

	return CmdResult{}
}

func KillUI() CmdResult {
	cmd := exec.Command("pidof", clientProcessName)
	out, err := cmd.Output()
	if err != nil {
		return CmdResult{err: err}
	}

	pid := bytes.TrimSpace(out)
	cmd = exec.Command("kill", string(pid))
	if err = cmd.Run(); err != nil {
		return CmdResult{err: err}
	}

	return CmdResult{}
}

func RestartServer() CmdResult {
	cmd := exec.Command("service", strings.Split(restartServer, " ")...)
	if err := cmd.Run(); err != nil {
		return CmdResult{err: err}
	}

	return CmdResult{}
}

func StartServer() CmdResult {
	cmd := exec.Command("service", strings.Split(startServer, " ")...)
	if err := cmd.Run(); err != nil {
		return CmdResult{err: err}
	}

	return CmdResult{}
}

func StopServer() CmdResult {
	cmd := exec.Command("service", strings.Split(stopServer, " ")...)
	if err := cmd.Run(); err != nil {
		return CmdResult{err: err}
	}

	return CmdResult{}
}

func SwitchToDebug() CmdResult {
	file, err := os.OpenFile(serverConfFilePath, os.O_RDWR, 0)
	if err != nil {
		return CmdResult{err: err}
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return CmdResult{err: err}
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

	return CmdResult{}
}
