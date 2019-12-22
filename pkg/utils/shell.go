package utils

import (
	"io"
	"os/exec"
	"path"
	"syscall"
)

type (
	// Executor execute task on local env
	Executor struct{}
)

// Exec - run a shell command
func (_e *Executor) Exec(command string, args []string, environ []string, detached bool, workingDir string, binpath string, logger io.Writer) (int, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = environ
	cmd.Stdout = logger
	cmd.Stderr = logger
	cmd.Dir = workingDir
	if binpath != "" {
		cmd.Path = path.Join(binpath, command)
	}
	if !detached {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		return cmd.Process.Pid, cmd.Run()
	} else {
		return cmd.Process.Pid, cmd.Start()
	}
}
