package utils

import (
	"fmt"
	"os/exec"
	"syscall"
)

type (
	Command struct {
		env     []string
		command []string
		bin     string
	}
)

func (_c *Command) AddEnv(key string, value string) {
	_c.env = append(_c.env, fmt.Sprintf("%s=%s", key, value))
}

func (_c *Command) AddCommand(cmd string) {
	_c.command = append(_c.command, cmd)
}

func (_c *Command) Bin(path string) {
	_c.bin = path
}

// Create builds exec.Cmd
func (_c *Command) Create() *exec.Cmd {
	cmd := exec.Command(_c.bin, _c.command...)
	cmd.Env = _c.env
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}
