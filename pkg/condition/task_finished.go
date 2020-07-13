package condition

import (
	"github.com/open-integration/core/pkg/event"
	"github.com/open-integration/core/pkg/state"
)

type conditionTaskFinished struct {
	task string
}

func (c *conditionTaskFinished) Met(ev event.Event, s state.State) bool {
	if ev.Metadata.Name != state.EventTaskFinished {
		return false
	}
	for _, t := range s.Tasks() {
		if t.State == state.TaskStateFinished && t.Task.Metadata().Name == c.task {
			return true
		}
	}
	return false
}
