package core

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/open-integration/core/internal/commands"
	"github.com/open-integration/core/pkg/logger"
	"github.com/open-integration/core/pkg/shell"
	"github.com/open-integration/core/pkg/utils"
)

type (

	// Engine exposes the interface of the engine
	Engine interface {
		Run() error
	}

	// EngineOptions to create new engine
	EngineOptions struct {
		Logger    logger.Logger
		Pipeline  Pipeline
		State     *State
		EventChan chan *Event
	}

	engine struct {
		pipeline         Pipeline
		logger           logger.Logger
		state            *State
		eventChan        chan *Event
		taskLogsDirctory string
		modem            Modem
	}
)

// Run starts the pipeline execution
func (e *engine) Run() error {
	e.logger.Debug("Starting...", "pipeline", e.pipeline.Metadata.Name)
	err := e.modem.Init()
	if err != nil {
		return err
	}
	defer e.modem.Destory()
	e.state.send(commands.StartEngine, func() *State {
		return &State{
			Metadata: StateMetadata{
				State: EngineStateInProgress,
			},
		}
	})
	e.handleStateEvents()
	return nil
}

func (e *engine) handleStateEvents() {
	for {
		ev := <-e.eventChan
		if ev.Metadata.Name == EventEngineFinished {
			e.logger.Debug("Received finish event, killing all services")
			return
		}
		go e.handleEvent(*ev)
	}
}

func (e *engine) handleEvent(ev Event) {
	e.logger.Debug("Recieved event", "name", ev.Metadata.Name)
	for _, t := range e.pipeline.Spec.Tasks {
		taskLogger := e.logger.New("task", t.Metadata.Name)
		if !e.shouldRunTask(t, &ev, taskLogger) {
			continue
		}
		e.runTask(t, &ev, taskLogger)
	}
	e.reportOnFinish()
}

func (e *engine) shouldRunTask(t Task, ev *Event, logger logger.Logger) bool {
	if !t.Metadata.Reusable {
		for _, tt := range e.state.Tasks {
			if tt.Task == t.Metadata.Name {
				// Skip tasks that are nnot resuable
				return false
			}
		}
	}

	logger.Debug("Running task conditions")
	if t.Condition == nil {
		logger.Debug("No condition set, skipping...")
		return false
	} else {
		logger.Debug("Running condition", "condition", t.Condition.Name)
		if !t.Condition.Func(ev, e.state) {
			logger.Debug("Condition evaludated to false, skipping...")
			return false
		}
	}
	return true
}

func (e *engine) runTask(t Task, ev *Event, logger logger.Logger) error {
	var spec TaskSpec

	if t.SpecFunc != nil {
		s, err := t.SpecFunc(e.state)
		if err != nil {
			return err
		}
		spec = *s
	} else {
		spec = t.Spec
	}
	id := generateID()
	e.state.send(commands.StartTask, func() *State {
		return &State{
			Tasks: map[ID]TaskState{
				id: TaskState{
					State:   TaskStateInProgress,
					Task:    t.Metadata.Name,
					EventID: ev.Metadata.ID,
				},
			},
		}
	})

	fileDescriptor := path.Join(e.taskLogsDirctory, fmt.Sprintf("%s-%s.log", t.Metadata.Name, string(id)))
	tl, err := utils.CreateLogFile(e.taskLogsDirctory, fmt.Sprintf("%s-%s.log", t.Metadata.Name, string(id)))
	if err != nil {
		logger.Error("Failed to create log file for task")
	}

	payload := ""
	if spec.Service != "" {
		e.logger.Debug("Calling service", "service", spec.Service, "endpoint", spec.Endpoint)
		payload, err = e.modem.Call(spec, fileDescriptor)
	} else {
		cmd := strings.Split(spec.Command, " ")
		_, err = shell.Execute(cmd[0], cmd[1:], append(spec.EnvironmentVariables, os.Environ()...), spec.Detached, spec.WorkingDirectory, spec.Path, tl)
	}

	status := TaskStatusSuccess
	msg := ""
	if err != nil {
		logger.Error("Task exited with error", "err", err.Error())
		status = TaskStatusFailed
		msg = err.Error()
	}
	e.state.send(commands.FinishTask, func() *State {
		return &State{
			Tasks: map[ID]TaskState{
				id: TaskState{
					State:   TaskStateFinished,
					Status:  status,
					Task:    t.Metadata.Name,
					Output:  payload,
					EventID: ev.Metadata.ID,
					Error:   msg,
				},
			},
		}
	})
	return nil
}

// reportOnFinish watch all events and send finish command once there are no more tasks in in-progress state
func (e *engine) reportOnFinish() {
	time.Sleep(time.Second * 3)

	for id, t := range e.state.Tasks {
		e.logger.Debug("Testing task", "name", t.Task, "id", id, "state", t.State)
		if t.State != "finished" {
			return
		}
	}

	e.logger.Debug("All tasks seems to be finished, sending finish command")
	go e.state.send(commands.FinishEngine, func() *State {
		return &State{
			Metadata: StateMetadata{
				State: EngineStateFinished,
			},
		}
	})
	return

}
