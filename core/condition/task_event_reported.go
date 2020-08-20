package condition

import (
	"fmt"

	"github.com/open-integration/core/core/event"
	"github.com/open-integration/core/core/state"
)

type conditionTaskEventReported struct {
	name string
}

func (c *conditionTaskEventReported) Met(ev event.Event, s state.State) bool {
	return fmt.Sprintf("%s.%s", ev.Metadata.Task, ev.Metadata.Name) == c.name
}
