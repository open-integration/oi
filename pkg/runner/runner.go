package runner

import (
	"github.com/open-integration/core/pkg/logger"
	"io"
	"os/exec"
)

const (
	// LocalRunner type
	LocalRunner = "local"

	// KubernetesRunner type
	KubernetesRunner = "kubernetes"
)

type (
	// Runner expose an interface to run services
	Runner interface {
		Run(log io.Writer) error
		Kill() error
	}

	// Options shows all the available options to build runner
	Options struct {
		Type           string
		Logger         logger.Logger
		LocalRunnerCmd *exec.Cmd
	}
)

// New builds new runner based on Options.Type
func New(opt *Options) Runner {
	if opt.Type == LocalRunner {
		return &localRunner{
			Logger:  opt.Logger,
			Command: opt.LocalRunnerCmd,
		}
	}

	if opt.Type == KubernetesRunner {
		return &kubernetesRunner{}
	}
	return nil
}
