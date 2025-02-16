package composer

import (
	"os"
	"path/filepath"
	"strings"
)

const FileName = "composer.json"

func ParseVersionConstraint(version string) string {
	if strings.HasPrefix(version, "^") || strings.HasPrefix(version, "~") {
		version = version[1:]
	}

	return version
}

func GetFilePath(path string) (string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if stat.IsDir() {
		path = filepath.Join(path, FileName)
	}

	return path, nil
}
