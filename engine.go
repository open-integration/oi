package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/open-integration/core/pkg/graph"
	"github.com/open-integration/core/pkg/logger"
	"github.com/open-integration/core/pkg/modem"
	"github.com/open-integration/core/pkg/state"
	"github.com/open-integration/core/pkg/task"
	"github.com/open-integration/core/pkg/utils"
)

type (

	// Engine exposes the interface of the engine
	Engine interface {
		Run() error
		Modem() modem.Modem
	}

	engine struct {
		pipeline           Pipeline
		logger             logger.Logger
		eventChan          chan *state.Event
		stateUpdateRequest chan state.StateUpdateRequest
		taskLogsDirctory   string
		stateDir           string
		modem              modem.Modem
		statev1            state.State
		wg                 *sync.WaitGroup
		graphBuilder       graph.Builder
	}
)

var now = func() time.Time { return time.Now() }

// Run starts the pipeline execution
func (e *engine) Run() error {
	e.logger.Debug("Starting...", "pipeline", e.pipeline.Metadata.Name)
	err := e.modem.Init()
	if err != nil {
		return err
	}
	defer func() {
		e.logger.Debug("killing all services")
		e.modem.Destroy()
	}()
	e.wg.Add(1)
	e.stateUpdateRequest <- state.StateUpdateRequest{
		Metadata: state.StateUpdateRequestMetadata{
			CreatedAt: now(),
		},
		UpdateStateMetadataRequest: &state.UpdateStateMetadataRequest{
			State: state.EngineStateInProgress,
		},
	}
	go e.waitForFinish()
	e.handleStateEvents()
	s, _ := e.statev1.Copy()
	g := e.graphBuilder.Build(s)
	ioutil.WriteFile(path.Join(e.stateDir, "graph.dot"), g, os.ModePerm)
	return e.printStateStore()
}

// Modem returns the current modem
func (e *engine) Modem() modem.Modem {
	return e.modem
}

// handleStateEvents watch the event channel and act on each evnt
// state.EventEngineFinished - finished watching, execution finished
// state.EventEngineStarted OR state.EventTaskStarted OR state.EventTaskFinished - elect next tasks
// state.EventTaskElected - execute tasks
func (e *engine) handleStateEvents() {
	for {
		ev := <-e.eventChan
		switch ev.Metadata.Name {
		case state.EventEngineFinished:
			return
		case state.EventEngineStarted, state.EventTaskStarted, state.EventTaskFinished:
			go e.electNextTasks(*ev)
		case state.EventTaskElected:
			go e.executeElectedTasks(*ev)
		}
	}
}

// electNextTasks - running all reactions on the event and sending request to elect matched tasks
func (e *engine) electNextTasks(ev state.Event) {
	log := e.logger.New("event", ev.Metadata.Name)
	log.Debug("Received event, electing next tasks")
	stateCpy, err := e.statev1.Copy()
	if err != nil {
		e.logger.Error("Failed to copy state")
		return
	}

	// candidate is the tasks results of all reactions
	tasksCandidates := map[string]task.Task{}
	for _, reaction := range e.pipeline.Spec.Reactions {
		if reaction.Condition(ev, stateCpy) {
			for _, t := range reaction.Reaction(ev, stateCpy) {
				tasksCandidates[t.Metadata.Name] = t
			}
		}
	}

	tasksToElect := []task.Task{}
	for _, t := range tasksCandidates {
		_, exist := stateCpy.Tasks()[t.Metadata.Name]
		if !exist {
			e.logger.Debug("Adding task to elected set", "task", t.Metadata.Name)
			tasksToElect = append(tasksToElect, t)
		}
	}
	if len(tasksToElect) > 0 {
		e.logger.Debug("Electing tasks", "total", len(tasksToElect))
		ids := []string{}
		for _, t := range tasksToElect {
			ids = append(ids, t.Metadata.Name)
		}
		e.stateUpdateRequest <- state.StateUpdateRequest{
			Metadata: state.StateUpdateRequestMetadata{
				CreatedAt: now(),
			},
			ElectTasksRequest: &state.ElectTasksRequest{
				Tasks: tasksToElect,
			},
			AddRealtedTaskToEventReuqest: &state.AddRealtedTaskToEventReuqest{
				EventID: ev.Metadata.ID,
				Task:    ids,
			},
		}
	}
}

// executeElectedTasks - execute all elected tasks in parallel
func (e *engine) executeElectedTasks(ev state.Event) {
	log := e.logger.New("event", ev.Metadata.Name)
	stateCpy, err := e.statev1.Copy()
	if err != nil {
		e.logger.Error("Failed to copy state")
		return
	}
	elected := []task.Task{}
	for _, t := range stateCpy.Tasks() {
		if t.State == state.TaskStateElected {
			elected = append(elected, t.Task)
		}
	}
	wg := &sync.WaitGroup{}
	for _, t := range elected {
		wg.Add(1)
		log.Debug("Running task", "task", t.Metadata.Name)
		go e.runTask(t, ev, log.New("task", t.Metadata.Name))
		wg.Done()
	}
	wg.Wait()
}

func (e *engine) runTask(t task.Task, ev state.Event, logger logger.Logger) {
	spec := t.Spec

	e.stateUpdateRequest <- state.StateUpdateRequest{
		Metadata: state.StateUpdateRequestMetadata{
			CreatedAt: now(),
		},
		AddRealtedTaskToEventReuqest: &state.AddRealtedTaskToEventReuqest{
			EventID: ev.Metadata.ID,
		},
	}
	fileName := fmt.Sprintf("%s.log", t.Metadata.Name)
	fileDescriptor := path.Join(e.taskLogsDirctory, fileName)
	e.wg.Add(1)
	e.stateUpdateRequest <- state.StateUpdateRequest{
		Metadata: state.StateUpdateRequestMetadata{
			CreatedAt: now(),
		},
		UpdateTaskStateRequest: &state.UpdateTaskStateRequest{
			State: state.TaskState{
				State:  state.TaskStateInProgress,
				Task:   t,
				Logger: fileDescriptor,
			},
		},
	}

	_, err := utils.CreateLogFile(e.taskLogsDirctory, fileName)
	if err != nil {
		logger.Error("Failed to create log file for task")
	}

	payload := ""
	arguments := map[string]interface{}{}
	if spec.Service != "" {
		e.logger.Debug("Calling service", "service", spec.Service, "endpoint", spec.Endpoint)
		for _, arg := range spec.Arguments {
			arguments[arg.Key] = arg.Value
		}
		payload, err = e.modem.Call(spec.Service, spec.Endpoint, arguments, fileDescriptor)
	}

	status := state.TaskStatusSuccess
	msg := ""
	if err != nil {
		logger.Error("Task exited with error", "err", err.Error())
		status = state.TaskStatusFailed
		msg = err.Error()
	}
	e.wg.Add(1)
	e.stateUpdateRequest <- state.StateUpdateRequest{
		Metadata: state.StateUpdateRequestMetadata{
			CreatedAt: now(),
		},
		UpdateTaskStateRequest: &state.UpdateTaskStateRequest{
			State: state.TaskState{
				State:  state.TaskStateFinished,
				Status: status,
				Task:   t,
				Output: payload,
				Error:  msg,
				Logger: fileDescriptor,
			},
		},
	}

}

// waitForFinish watch all events and send finish command once there are no more tasks in in-progress state
func (e *engine) waitForFinish() {
	time.Sleep(5 * time.Second)
	e.wg.Wait()
	stateCpy, _ := e.statev1.Copy()
	for _, t := range stateCpy.Tasks() {
		if t.State != "finished" {
			go e.waitForFinish()
			return
		}
	}

	e.logger.Debug("All tasks seems to be finished, sending finish command")
	e.wg.Add(1)
	e.stateUpdateRequest <- state.StateUpdateRequest{
		Metadata: state.StateUpdateRequestMetadata{
			CreatedAt: now(),
		},
		UpdateStateMetadataRequest: &state.UpdateStateMetadataRequest{
			State:  state.EngineStateFinished,
			Status: state.EngineStatusSuccess,
		},
	}
	return
}

func (e *engine) printStateStore() error {
	statebytes, err := e.statev1.StateBytes()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(e.stateDir, "state.yaml"), statebytes, os.ModePerm)
	if err != nil {
		e.logger.Error("Failed to store state to file")
		return err
	}

	eventbytes, err := e.statev1.EventBytes()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(e.stateDir, "events.yaml"), eventbytes, os.ModePerm)
	if err != nil {
		e.logger.Error("Failed to store state to file")
		return err
	}
	return nil
}
