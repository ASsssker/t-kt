package commands

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
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

func DisableObselete() CmdResult {
	repFiles := []string{}
	f, err := os.Open(DriverPackRepoAbsPath)
	if err != nil {
		return CmdResult{err: err}
	}
	defer f.Close()

	files, err := f.ReadDir(0)
	if err != nil {
		return CmdResult{err: err}
	}

	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), ".rep") {
			repFiles = append(repFiles, DriverPackRepoAbsPath+"/"+file.Name())
		}
	}

	errChan := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(len(repFiles))
	for _, fileName := range repFiles {
		go func(fileName string) {
			defer wg.Done()
			file, err := os.OpenFile(fileName, os.O_RDWR, 0)
			if err != nil {
				errChan <- err
				return
			}
			defer file.Close()

			buf, err := io.ReadAll(file)
			if err != nil {
				errChan <- err
				return
			}

			buf = bytes.Replace(buf, []byte("obsolete=\"true\""), []byte("obsolete=\"false\""), -1)

			if err = file.Truncate(0); err != nil {
				errChan <- err
				return
			}
			if _, err = file.Seek(0, 0); err != nil {
				errChan <- err
				return
			}

			if _, err = file.Write(buf); err != nil {
				errChan <- err
				return
			}
		}(fileName)
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return CmdResult{}
	case err := <-errChan:
		return CmdResult{err: err}
	}
}
