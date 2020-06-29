package condition

import (
	"github.com/open-integration/core/pkg/event"
	"github.com/open-integration/core/pkg/state"
)

type conditionEngineStarted struct{}

func (c *conditionEngineStarted) Met(ev event.Event, s state.State) bool {
	return ev.Metadata.Name == state.EventEngineStarted
}
