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
