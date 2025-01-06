package composer

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type File struct {
	Require map[string]string `json:"require"`
}

func (f *File) GetPhpVersion() (string, error) {
	version, ok := f.Require["php"]
	if !ok {
		return "", errors.New("php version not found")
	}

	return ParseVersionConstraint(version), nil
}

func Parse(path string) (*File, error) {
	path, err := GetFilePath(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var composer File
	if err := json.Unmarshal(data, &composer); err != nil {
		return nil, err
	}

	return &composer, nil
}
