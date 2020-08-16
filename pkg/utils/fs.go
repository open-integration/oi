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
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("Failed to create dir %s: %w", dir, err)
	}
	path := fmt.Sprintf("%s/%s", dir, name)
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to create file %s: %w", name, err)
	}
	return f, nil
}
