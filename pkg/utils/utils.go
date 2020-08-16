package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/gofrs/uuid"
)

// CreateLogFile creates file to be used as logger
func CreateLogFile(dir string, name string) (io.Writer, error) {
	if dir == "" {
		return os.Stdout, nil
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("Failed to create logs directory in %s: %w", dir, err)
	}
	path := fmt.Sprintf("%s/%s", dir, name)
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to create log file %s: %w", path, err)
	}
	return f, nil
}

// GenerateID build random uuid-v4
func GenerateID() string {
	return uuid.Must(uuid.NewV4()).String()
}
