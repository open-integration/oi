package task

import (
	"context"
	"time"

	"github.com/open-integration/core/pkg/event/reporter"
	"github.com/open-integration/core/pkg/filedescriptor"
	"github.com/open-integration/core/pkg/modem"
)

type (
	Task interface {
		Run(ctx context.Context, options RunOptions) ([]byte, error)
		Metadata() Metadata
	}

	RunOptions struct {
		FD            filedescriptor.FD
		EventReporter reporter.EventReporter
		Modem         modem.ServiceCaller
	}

	// Metadata holds all the metadata of a pipeline
	Metadata struct {
		Name string
		Time struct {
			// StaredAt the time when the step was started execution
			StartedAt time.Time
			// FinishedAt the time when the step finished execution
			FinishedAt time.Time
		}
	}

	// Argument is key value struct that should be passed in a service call
	Argument struct {
		Key   string
		Value interface{}
	}
)
