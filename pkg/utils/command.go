package utils

import (
	"fmt"
	"os/exec"
	"syscall"
)

type (
	// Command runs command in as subprocess
	Command struct {
		env     []string
		command []string
		bin     string
	}
)

// AddEnv adds env variable to the command
func (_c *Command) AddEnv(key string, value string) {
	_c.env = append(_c.env, fmt.Sprintf("%s=%s", key, value))
}

// AddCommand adds command argument to the command
func (_c *Command) AddCommand(cmd string) {
	_c.command = append(_c.command, cmd)
}

// Bin sets the binary to execute
func (_c *Command) Bin(path string) {
	_c.bin = path
}

// Create builds exec.Cmd
func (_c *Command) Create() *exec.Cmd {
	// #nosec
	cmd := exec.Command(_c.bin, _c.command...)
	cmd.Env = _c.env
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}
