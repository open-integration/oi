package condition

import (
	"github.com/open-integration/open-integration/core/event"
	"github.com/open-integration/open-integration/core/state"
)

type conditionCombined struct {
	conditions []Condition
}

func (c *conditionCombined) Met(ev event.Event, s state.State) bool {
	result := true
	for _, c := range c.conditions {
		if met := c.Met(ev, s); !met {
			result = false
		}
	}
	return result
}
