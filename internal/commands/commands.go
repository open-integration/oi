package commands

const (
	// StartEngine command will be sent once
	StartEngine Command = iota
	// FinishEngine command will be sent once at the end of the execution
	FinishEngine

	// StartTask
	StartTask

	// FinishTask
	FinishTask

	startEngine  = "engine:start"
	finishEngine = "engine:finish"
	startTask    = "task:start"
	finishTask   = "task:finish"
)

type (
	// Command enums
	Command int
)

func (n Command) String() string {
	return []string{
		startEngine,
		finishEngine,
		startTask,
		finishTask,
	}[n]
}
