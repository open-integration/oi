package condition

import "github.com/open-integration/core/pkg/state"

type conditionCombined struct {
	conditions []Condition
}

func (c *conditionCombined) Met(ev state.Event, s state.State) bool {
	result := true
	for _, c := range c.conditions {
		if met := c.Met(ev, s); !met {
			result = false
		}
	}
	return result
}
