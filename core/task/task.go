package task

import (
	"context"

	"github.com/open-integration/oi/core/event/reporter"
	"github.com/open-integration/oi/core/filedescriptor"
	"github.com/open-integration/oi/core/modem"
)

type (
	// Task defines task interface
	Task interface {
		Run(ctx context.Context, options RunOptions) ([]byte, error)
		Name() string
	}

	// RunOptions run tasks with given options
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
