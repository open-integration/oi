package condition

import (
	"github.com/open-integration/oi/core/event"
	"github.com/open-integration/oi/core/state"
)

type conditionTaskFinished struct {
	task string
}

func (c *conditionTaskFinished) Met(ev event.Event, s state.State) bool {
	if ev.Metadata.Name != state.EventTaskFinished {
		return false
	}
	for _, t := range s.Tasks() {
		if t.State == state.TaskStateFinished && t.Task.Name() == c.task {
			return true
		}
	}
	return false
}
