package runner

import (
	"context"

	v1 "github.com/open-integration/core/pkg/api/v1"
)

type (
	kubernetesRunner struct{}
)

func (_k *kubernetesRunner) Run() error {
	return nil
}

func (_k *kubernetesRunner) Kill() error {
	return nil
}
func (_k *kubernetesRunner) Call(context context.Context, req *v1.CallRequest) (*v1.CallResponse, error) {
	return nil, nil
}
