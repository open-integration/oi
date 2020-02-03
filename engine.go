package core

import (
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/open-integration/core/internal/commands"
	"github.com/open-integration/core/pkg/logger"
	"github.com/open-integration/core/pkg/modem"
	"github.com/open-integration/core/pkg/utils"
)

type (

	// Engine exposes the interface of the engine
	Engine interface {
		Run() error
	}

	// EngineOptions to create new engine
	EngineOptions struct {
		Pipeline   Pipeline
		Kubeconfig *EngineKubernetesOptions
	}

	EngineKubernetesOptions struct {
		Path      string
		Context   string
		Namespace string
		InCluster bool
	}

	engine struct {
		pipeline         Pipeline
		logger           logger.Logger
		state            *State
		eventChan        chan *Event
		taskLogsDirctory string
		modem            modem.Modem
		wg               sync.WaitGroup
	}
)

// Run starts the pipeline execution
func (e *engine) Run() error {
	e.logger.Debug("Starting...", "pipeline", e.pipeline.Metadata.Name)
	err := e.modem.Init()
	if err != nil {
		return err
	}
	defer e.modem.Destroy()
	e.wg.Add(1)
	go e.state.send(commands.StartEngine, &e.wg, func() *State {
		return &State{
			Metadata: StateMetadata{
				State: EngineStateInProgress,
			},
		}
	})
	go e.waitForFinish()
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
	e.wg.Add(1)
	log := e.logger.New("event", ev.Metadata.Name)
	log.Debug("Received event", "total-reactions", len(e.pipeline.Spec.Reactions))
	for _, reaction := range e.pipeline.Spec.Reactions {
		log.Debug("Running reaction condition")
		if reaction.Condition(ev, *e.state) {
			log.Debug("Condition evaluated to true")
			tasks := reaction.Reaction(ev, *e.state)
			log.Debug("Received new tasks", "len", len(tasks))
			for _, t := range tasks {
				taskLogger := log.New("task", t.Metadata.Name)
				if !t.Metadata.Reusable {
					shouldSkip := false
					for _, pastTask := range e.state.Tasks {
						if pastTask.Task == t.Metadata.Name {
							taskLogger.Debug("Task been executed in the past, skiping")
							shouldSkip = true
						}
					}
					if shouldSkip {
						continue
					}
				}
				err := e.runTask(t, &ev, taskLogger)
				if err != nil {
					log.Error("Error running task", "err", err.Error(), "task", t.Metadata.Name)
				}
			}
		} else {
			log.Debug("Reaction condition evaludated to false")
		}
	}
	e.wg.Done()
}

func (e *engine) runTask(t Task, ev *Event, logger logger.Logger) error {
	spec := t.Spec
	id := generateID()
	e.wg.Add(1)
	go e.state.send(commands.StartTask, &e.wg, func() *State {
		return &State{
			Tasks: map[ID]TaskState{
				id: {
					State:   TaskStateInProgress,
					Task:    t.Metadata.Name,
					EventID: ev.Metadata.ID,
				},
			},
		}
	})

	fileName := fmt.Sprintf("%s-%s.log", t.Metadata.Name, string(id))
	fileDescriptor := path.Join(e.taskLogsDirctory, fileName)
	_, err := utils.CreateLogFile(e.taskLogsDirctory, fileName)
	if err != nil {
		logger.Error("Failed to create log file for task")
	}

	payload := ""
	if spec.Service != "" {
		e.logger.Debug("Calling service", "service", spec.Service, "endpoint", spec.Endpoint)
		arguments := map[string]interface{}{}
		for _, arg := range spec.Arguments {
			arguments[arg.GetKey()] = arg.GetValue()
		}
		payload, err = e.modem.Call(spec.Service, spec.Endpoint, arguments, fileDescriptor)
	}

	status := TaskStatusSuccess
	msg := ""
	if err != nil {
		logger.Error("Task exited with error", "err", err.Error())
		status = TaskStatusFailed
		msg = err.Error()
	}
	e.wg.Add(1)
	go e.state.send(commands.FinishTask, &e.wg, func() *State {
		return &State{
			Tasks: map[ID]TaskState{
				id: {
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

// waitForFinish watch all events and send finish command once there are no more tasks in in-progress state
func (e *engine) waitForFinish() {
	time.Sleep(5 * time.Second)
	e.wg.Wait()

	for id, t := range e.state.Tasks {
		e.logger.Debug("Testing task", "name", t.Task, "id", id, "state", t.State)
		if t.State != "finished" {
			return
		}
	}

	e.logger.Debug("All tasks seems to be finished, sending finish command")
	e.wg.Add(1)
	go e.state.send(commands.FinishEngine, &e.wg, func() *State {
		return &State{
			Metadata: StateMetadata{
				State: EngineStateFinished,
			},
		}
	})
	return

}
