package utils

import (
	"net"
)

type (
	Port struct{}
)

// GetAvailable returns available port on local machine
func (_p Port) GetAvailable() (string, error) {
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
