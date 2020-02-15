package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/gofrs/uuid"
)

func CreateLogFile(dir string, name string) (io.Writer, error) {
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

func GenerateID() string {
	return uuid.Must(uuid.NewV4()).String()
}
