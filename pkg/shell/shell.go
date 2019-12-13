package shell

import (
	"io"
	"os/exec"
	"path"
	"syscall"
)

// Execute - run a shell command
func Execute(command string, args []string, environ []string, detached bool, workingDir string, binpath string, logger io.Writer) (int, error) {
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
