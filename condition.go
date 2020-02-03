package core

// ConditionEngineStarted returns true once the engine emits started event
func ConditionEngineStarted(ev Event, state State) bool {
	return ev.Metadata.Name == EventEngineStarted
}

// ConditionTaskFinishedWithStatus returns true once a task is finished with given status
func ConditionTaskFinishedWithStatus(task string, status string) func(ev Event, state State) bool {
	return func(ev Event, state State) bool {
		if ev.Metadata.Name != EventTaskFinished {
			return false
		}
		for _, t := range state.Tasks {
			if t.Status == status && t.State == "finished" && t.Task == task {
				return true
			}
		}
		return false
	}
}
