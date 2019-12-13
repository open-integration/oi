package utils

import (
	"fmt"
	"io"
	"net"
	"os"
)

func GetAvailablePort() (string, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}
	defer l.Close()

	_, p, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return "", err
	}
	return p, err
}

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
