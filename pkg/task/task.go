package task

import (
	"context"

	"github.com/open-integration/core/pkg/event/reporter"
	"github.com/open-integration/core/pkg/filedescriptor"
	"github.com/open-integration/core/pkg/modem"
)

type (
	Task interface {
		Run(ctx context.Context, options RunOptions) ([]byte, error)
		Name() string
	}

	RunOptions struct {
		FD            filedescriptor.FD
		EventReporter reporter.EventReporter
		Modem         modem.ServiceCaller
	}

	// Metadata holds all the metadata of a pipeline
	Metadata struct {
		Name string
	}

	// Argument is key value struct that should be passed in a service call
	Argument struct {
		Key   string
		Value interface{}
	}
)
