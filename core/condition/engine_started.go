package condition

import (
	"github.com/open-integration/open-integration/core/event"
	"github.com/open-integration/open-integration/core/state"
)

type conditionEngineStarted struct{}

func (c *conditionEngineStarted) Met(ev event.Event, s state.State) bool {
	return ev.Metadata.Name == state.EventEngineStarted
}
