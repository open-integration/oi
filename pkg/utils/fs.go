package utils

import (
	"fmt"
	"io"
	"os"
)

type (
	// FileCreator expose ability to operate on fs
	FileCreator struct{}
)

// Create creates file in directory with mode = os.ModePerm
func (_f *FileCreator) Create(dir string, name string) (io.Writer, error) {
	if dir == "" {
		return os.Stdout, nil
	}
	os.MkdirAll(dir, os.ModePerm)
	path := fmt.Sprintf("%s/%s", dir, name)
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
