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
		return nil, fmt.Errorf("failed to create logs directory in %s: %w", dir, err)
	}
	path := fmt.Sprintf("%s/%s", dir, name)
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file %s: %w", path, err)
	}
	return f, nil
}

// GenerateID build random uuid-v4
func GenerateID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func GetEnvOrDefault(name string, def string) string {
	res := os.Getenv(name)
	if res == "" {
		return def
	}
	return res
}

func GetEnvOrDie(name string, errMsg string) string {
	if res := os.Getenv(name); res != "" {
		return res
	}
	DieOnError(fmt.Errorf("%s is not set", name), errMsg)
	return ""
}

func DieOnError(err error, msg string) {
	if err == nil {
		return
	}
	fmt.Printf("Error: %s\n", msg)
	panic(err)
}
