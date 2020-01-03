package runner

import (
	"io"
	"os"
	"os/exec"

	"github.com/open-integration/core/pkg/logger"
)

type (
	localRunner struct {
		Command *exec.Cmd
		Logger  logger.Logger
	}
)

func (_l *localRunner) Run(log io.Writer) error {
	_l.Logger.Debug("Starting service")
	_l.Command.Stdout = log
	_l.Command.Stderr = log
	if _l.Command.SysProcAttr.Setpgid {
		return _l.Command.Start()
	}
	return _l.Command.Run()
}

func (_l *localRunner) Kill() error {
	_l.Logger.Debug("Killing service")
	process, err := os.FindProcess(_l.Command.Process.Pid)
	if err != nil {
		return err
	}
	return process.Signal(os.Interrupt)
}
