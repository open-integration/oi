package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/open-integration/core/pkg/logger"
	"github.com/open-integration/core/pkg/modem"
	"github.com/open-integration/core/pkg/state"
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
		modem              modem.Modem
		statev1            state.State
		wg                 *sync.WaitGroup
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
	return e.printStateStore()
}

// Modem returns the current modem
func (e *engine) Modem() modem.Modem {
	return e.modem
}

func (e *engine) handleStateEvents() {
	for {
		ev := <-e.eventChan
		if ev.Metadata.Name == state.EventEngineFinished {
			return
		}
		go e.handleEvent(*ev)
	}
}

// handleEvent
// Run over all event reactions and get the task candidates if the condition is passed
// Run over all the candidates and get the task that actually should be executed
// based on the .metadata.reusable prop
func (e *engine) handleEvent(ev state.Event) {

	log := e.logger.New("event", ev.Metadata.Name)
	log.Debug("Received event")
	stateCpy, err := e.statev1.Copy()
	if err != nil {
		e.logger.Error("Failed to copy state")
		return
	}

	// candidate is the tasks results of all reactions
	tasksCandidates := []Task{}
	for _, reaction := range e.pipeline.Spec.Reactions {
		if reaction.Condition(ev, stateCpy) {
			tasksCandidates = append(tasksCandidates, reaction.Reaction(ev, stateCpy)...)
		}
	}

	pairs := []struct {
		task   Task
		logger logger.Logger
	}{}
	tasksToElect := []string{}
	for _, t := range tasksCandidates {
		taskLogger := log.New("task", t.Metadata.Name)
		_, exist := stateCpy.Tasks()[t.Metadata.Name]
		if !exist {
			pairs = append(pairs, struct {
				task   Task
				logger logger.Logger
			}{
				task:   t,
				logger: taskLogger,
			})
			tasksToElect = append(tasksToElect, t.Metadata.Name)
		}

	}
	e.logger.Debug("Electing tasks", "total", len(tasksToElect))
	e.electTasks(tasksToElect)
	for _, pair := range pairs {
		err := e.runTask(pair.task, ev, pair.logger)
		if err != nil {
			log.Error("Error running task", "err", err.Error(), "task", pair.task.Metadata.Name)
		}
	}
}

func (e *engine) electTasks(tasks []string) {
	e.stateUpdateRequest <- state.StateUpdateRequest{
		Metadata: state.StateUpdateRequestMetadata{
			CreatedAt: now(),
		},
		ElectTasksRequest: &state.ElectTasksRequest{
			Tasks: tasks,
		},
	}
}

func (e *engine) runTask(t Task, ev state.Event, logger logger.Logger) error {
	spec := t.Spec

	e.stateUpdateRequest <- state.StateUpdateRequest{
		Metadata: state.StateUpdateRequestMetadata{
			CreatedAt: now(),
		},
		AddRealtedTaskToEventReuqest: &state.AddRealtedTaskToEventReuqest{
			EventID: ev.Metadata.ID,
			Task:    t.Metadata.Name,
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
				State:   state.TaskStateInProgress,
				Task:    t.Metadata.Name,
				EventID: ev.Metadata.ID,
				Logger:  fileDescriptor,
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
			arguments[arg.GetKey()] = arg.GetValue()
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
				State:     state.TaskStateFinished,
				Status:    status,
				Task:      t.Metadata.Name,
				Output:    payload,
				Arguments: fmt.Sprintf("%v", arguments),
				EventID:   ev.Metadata.ID,
				Error:     msg,
				Logger:    fileDescriptor,
			},
		},
	}

	return nil
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
			State: state.EngineStateFinished,
		},
	}
	return
}

func (e *engine) printStateStore() error {
	statebytes, err := e.statev1.StateBytes()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./logs/state.yaml", statebytes, os.ModePerm)
	if err != nil {
		e.logger.Error("Failed to store state to file")
		return err
	}

	eventbytes, err := e.statev1.EventBytes()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./logs/events.yaml", eventbytes, os.ModePerm)
	if err != nil {
		e.logger.Error("Failed to store state to file")
		return err
	}
	return nil
}
