package background

import (
	"fmt"
	"os"
	"strings"
	"t-kt/internal/commands"
	"time"
)

var dumpList = make(map[string]struct{})
var checkDumpInit = false

func CheckDump() commands.CmdResult {
	clientLogDir, err := commands.GetClientLogDir()
	if err != nil {
		return commands.NewCmdResult("", err)
	}

	clientDumps, err := getDumps(clientLogDir)
	if err != nil {
		return commands.NewCmdResult("", err)
	}
	serverDumps, err := getDumps(commands.ServerLogAbsPath)
	if err != nil {
		return commands.NewCmdResult("", err)
	}

	allDumps := mergeMaps(clientDumps, serverDumps)

	newDump := findDump(allDumps, dumpList)
	if newDump != "" {
		dumpList = mergeMaps(allDumps, dumpList)
		msg := fmt.Sprintf("%s падение дампа:%s", time.Now().Format("15:04:05"), newDump)

		// при первом запуске игнорируем дампы
		if checkDumpInit {
			// ограничение на частоту запуска в секунду
			return commands.NewCmdResult(msg, nil)
		} else {
			checkDumpInit = true
		}
	}

	return commands.NewCmdResult("", nil)
}

func findDump(newDumps map[string]struct{}, oldDumps map[string]struct{}) string {
	for key := range newDumps {
		if _, exitsts := oldDumps[key]; !exitsts {
			return key
		}
	}

	return ""
}

func mergeMaps[T comparable](m1, m2 map[T]struct{}) map[T]struct{} {
	for k := range m1 {
		m2[k] = struct{}{}
	}

	return m2
}

func getDumps(path string) (map[string]struct{}, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	dumps := make(map[string]struct{})
	for _, file := range dir {
		if strings.Contains(file.Name(), "dmp") {
			dumps[fmt.Sprintf("%s/%s", path, file.Name())] = struct{}{}
		}
	}
	return dumps, nil
}
