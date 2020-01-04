package utils

import (
	"fmt"
	"os/exec"
	"syscall"
)

type (
	Command struct{}
)

// Create build exec.Cmd
func (_c Command) Create(port string, path string) *exec.Cmd {
	cmd := exec.Command(path)
	envs := []string{
		fmt.Sprintf("PORT=%s", port),
	}
	cmd.Env = envs
	cmd.Dir = ""
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}
