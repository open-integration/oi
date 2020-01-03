package runner

import "io"

type (
	kubernetesRunner struct{}
)

func (_k *kubernetesRunner) Run(log io.Writer) error {
	return nil
}

func (_k *kubernetesRunner) Kill() error {
	return nil
}
