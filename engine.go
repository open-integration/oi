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
		pipeline         Pipeline
		logger           logger.Logger
		eventChan        chan *state.Event
		taskLogsDirctory string
		modem            modem.Modem
		wg               sync.WaitGroup
		statev1          state.State
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
	ch := e.statev1.Send(state.CmdStartEngine, &e.wg)
	ch <- state.Change{
		Cause: "Pipeline",
		PipelineChange: state.PipelineChange{
			States: state.EngineStateInProgress,
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
			e.logger.Debug("Received finish event, killing all services")
			return
		}
		go e.handleEvent(*ev)
	}
}

// handleEvent
// Run over all event reactions and get the task candidates if the condition is passed
// Run over all the candidates and get the task that actually should be executed
// based on the .metadata.resuable prop
func (e *engine) handleEvent(ev state.Event) {
	e.wg.Add(1)
	log := e.logger.New("event", ev.Metadata.Name)
	stateCpy, _ := e.statev1.Copy()
	log.Debug("Received event", "total-reactions", len(e.pipeline.Spec.Reactions))
	tasksCandidates := []Task{}
	for _, reaction := range e.pipeline.Spec.Reactions {
		log.Debug("Running reaction condition")
		if reaction.Condition(ev, stateCpy) {
			log.Debug("Condition evaluated to true")
			tasksCandidates = append(tasksCandidates, reaction.Reaction(ev, stateCpy)...)
		} else {
			log.Debug("Reaction condition evaludated to false")
		}
	}

	paris := []struct {
		task   Task
		logger logger.Logger
	}{}
	for _, t := range tasksCandidates {
		taskLogger := log.New("task", t.Metadata.Name)
		if !t.Metadata.Reusable {
			shouldSkip := false
			for _, pastTask := range stateCpy.Tasks() {
				if pastTask.Task == t.Metadata.Name {
					taskLogger.Debug("Task been executed in the past, skiping")
					shouldSkip = true
				}
			}
			if !shouldSkip {
				paris = append(paris, struct {
					task   Task
					logger logger.Logger
				}{
					task:   t,
					logger: taskLogger,
				})
			}
		}
	}
	for _, pair := range paris {
		err := e.runTask(pair.task, &ev, pair.logger)
		if err != nil {
			log.Error("Error running task", "err", err.Error(), "task", pair.task.Metadata.Name)
		}
	}
	e.wg.Done()
}

func (e *engine) runTask(t Task, ev *state.Event, logger logger.Logger) error {
	spec := t.Spec
	id := generateID()
	e.wg.Add(1)
	fileName := fmt.Sprintf("%s-%s.log", t.Metadata.Name, string(id))
	fileDescriptor := path.Join(e.taskLogsDirctory, fileName)
	ch := e.statev1.Send(state.CmdStartTask, &e.wg)
	ch <- state.Change{
		Cause: "Task started",
		TaskChanges: state.TaskChanges{
			ID:       string(id),
			State:    state.TaskStateInProgress,
			Name:     t.Metadata.Name,
			EventID:  ev.Metadata.ID,
			LoggerID: fileDescriptor,
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
	ch = e.statev1.Send(state.CmdFinishTask, &e.wg)
	ch <- state.Change{
		Cause: "Task finished",
		TaskChanges: state.TaskChanges{
			ID:        string(id),
			State:     state.TaskStateFinished,
			Status:    status,
			Name:      t.Metadata.Name,
			Output:    payload,
			Arguments: fmt.Sprintf("%v", arguments),
			EventID:   ev.Metadata.ID,
			Error:     msg,
			LoggerID:  fileDescriptor,
		},
	}
	return nil
}

// waitForFinish watch all events and send finish command once there are no more tasks in in-progress state
func (e *engine) waitForFinish() {
	time.Sleep(5 * time.Second)
	e.wg.Wait()
	stateCpy, _ := e.statev1.Copy()
	for id, t := range stateCpy.Tasks() {
		e.logger.Debug("Testing task", "name", t.Task, "id", id, "state", t.State)
		if t.State != "finished" {
			return
		}
	}

	e.logger.Debug("All tasks seems to be finished, sending finish command")
	e.wg.Add(1)
	ch := e.statev1.Send(state.CmdFinishEngine, &e.wg)
	ch <- state.Change{
		Cause: "Pipeline finished",
		PipelineChange: state.PipelineChange{
			States: state.EngineStateFinished,
		},
	}
	return
}

func (e *engine) printStateStore() error {
	bytes, err := e.statev1.Bytes()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./logs/state.yaml", bytes, os.ModePerm)
	if err != nil {
		e.logger.Error("Failed to store state to file")
		return err
	}
	return nil
}
