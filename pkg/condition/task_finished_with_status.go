package condition

import "github.com/open-integration/core/pkg/state"

type conditionTaskFinishedWithStatus struct {
	task   string
	status string
}

func (c *conditionTaskFinishedWithStatus) Met(ev state.Event, s state.State) bool {
	if ev.Metadata.Name != state.EventTaskFinished {
		return false
	}
	for _, t := range s.Tasks() {
		if t.Status == c.status && t.State == state.TaskStateFinished && t.Task.Metadata.Name == c.task {
			return true
		}
	}
	return false
}
