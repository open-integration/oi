package reporter

import (
	"time"

	"github.com/open-integration/open-integration/core/event"
	"github.com/open-integration/open-integration/pkg/utils"
)

type (
	// EventReporter an interface that is used to report events during task execution
	EventReporter interface {
		Report(string, interface{}) error
	}

	// Options to create new EventReporter
	Options struct {
		EventChan chan<- event.Event
		Name      string
	}

	reporter struct {
		eventChan chan<- event.Event
		task      string
	}
)

// New creates new EventReporter based on Options
func New(opt Options) EventReporter {
	return reporter{
		eventChan: opt.EventChan,
		task:      opt.Name,
	}
}

func (r reporter) Report(name string, payload interface{}) error {
	r.eventChan <- event.Event{
		Metadata: event.Metadata{
			CreatedAt: time.Now(),
			ID:        utils.GenerateID(),
			Name:      name,
			Task:      r.task,
		},
		RelatedTasks: []string{},
		Payload: map[string]interface{}{
			"payload": payload,
		},
	}
	return nil
}
