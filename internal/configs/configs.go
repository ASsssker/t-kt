package configs

import (
	"errors"
	"strings"
)

type Config struct {
	IPC IPCConf
}


type IPCConf struct {
	DNS string
	ArchiveSwitchTime int
}

func (i *IPCConf) GetDNSConf() (username, password, addr string, err error) {
	temp := strings.Split(i.DNS, ":")
	if len(temp) <= 1 {
		return "", "", "", errors.New("invalid IPC data")
	}
	username = temp[0]
	temp = strings.Split(temp[1], "@")
	if len(temp) <= 1 {
		return "", "", "", errors.New("invalid IPC data")
	}
	password = temp[0]
	addr = temp[1]

	return username, password, addr, nil
}