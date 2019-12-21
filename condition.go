package core

type (
	// Condition determine whenever run a task
	Condition struct {
		// Name of the condition
		Name string
		// Func runs all the times
		Func func(*Event, *State) bool
	}
)

// ConditionEngineStarted returns true once the engine emits started event
func ConditionEngineStarted(ev *Event, state *State) bool {
	return ev.Metadata.Name == "engine.started"
}

// ConditionTaskFinishedWithStatus returns true once a task is finished with given status
func ConditionTaskFinishedWithStatus(task string, status string) func(ev *Event, state *State) bool {
	return func(ev *Event, state *State) bool {
		for _, t := range state.Tasks {
			if t.Status == status && t.State == "finished" && t.Task == task {
				return true
			}
		}
		return false
	}
}
