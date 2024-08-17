package commands

import (
	"fmt"
	"os"
	"path"
)

func clearDir(dirPath string) error {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range dir {
		err = os.RemoveAll(path.Join(dirPath, file.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

func getClientLogDir() (string, error) {
	user := os.Getenv("SUDO_USER")
	if user == "" {
		return "", fmt.Errorf("не удалось получить путь к клиентским логам")
	}

	return fmt.Sprintf("/home/%s%s", user, clientLogRelativePath), nil
}
