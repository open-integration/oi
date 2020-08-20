package tasks

import (
	"context"

	"github.com/open-integration/oi/core/task"
)

type (
	// svc calls a give service to run the task
	svc struct {
		meta      task.Metadata
		name      string
		endpoint  string
		arguments []task.Argument
	}
)

func (s *svc) Run(ctx context.Context, options task.RunOptions) ([]byte, error) {
	arguments := map[string]interface{}{}
	for _, arg := range s.arguments {
		arguments[arg.Key] = arg.Value
	}
	res, err := options.Modem.Call(ctx, s.name, s.endpoint, arguments, options.FD.File())
	return res, err
}

func (s *svc) Name() string {
	return s.meta.Name
}
