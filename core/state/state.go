package state

import (
	"errors"
	"sync"
	"time"

	"github.com/imdario/mergo"
	"github.com/open-integration/core/core/event"
	"github.com/open-integration/core/pkg/logger"
	"github.com/open-integration/core/pkg/utils"
	"gopkg.in/yaml.v2"
)

type (
	// State holds all the data of the pipeline execution flow
	State interface {
		Copy() (State, error)
		Tasks() map[string]TaskState
		Events() []event.Event
		Services() []ServiceState
		StateBytes() ([]byte, error)
		EventBytes() ([]byte, error)
		StartProcess()
	}

	state struct {
		name               string
		state              string
		services           []ServiceState
		tasks              map[string]TaskState
		events             []event.Event
		eventChan          chan *event.Event
		commandsChan       chan string
		logger             logger.Logger
		initialized        bool
		copyChan           chan state
		copyChanRequest    chan bool
		stateUpdateRequest chan UpdateRequest
		wg                 *sync.WaitGroup
	}

	// Options to pass to the state
	Options struct {
		Name string
		// EventChan to write new event to the channel once a chance was applied
		EventChan chan *event.Event
		// CommandsChan to receive commands to create new change channel
		CommandsChan chan string
		// StateUpdateRequest to receive updated on the state in realtime
		StateUpdateRequest chan UpdateRequest
		Logger             logger.Logger
		WG                 *sync.WaitGroup
	}
)

// New builds nnew State from options
func New(opt *Options) State {
	return &state{
		eventChan:          opt.EventChan,
		commandsChan:       opt.CommandsChan,
		logger:             opt.Logger,
		initialized:        false,
		copyChan:           make(chan state, 1),
		copyChanRequest:    make(chan bool, 1),
		stateUpdateRequest: opt.StateUpdateRequest,
		tasks:              map[string]TaskState{},
		wg:                 opt.WG,
	}
}

// Copy
func (s *state) Copy() (State, error) {
	s.copyChanRequest <- true
	for {
		select {
		case cpy, _ := <-s.copyChan:
			return &cpy, nil
		case <-time.After(60 * time.Second):
			msg := "Failed to copy state after 10 seconds"
			return nil, errors.New(msg)
		}
	}
}
func (s *state) StateBytes() ([]byte, error) {
	res := map[string]interface{}{
		"metadata": map[string]interface{}{
			"state": s.state,
		},
		"tasks":   s.tasks,
		"service": s.services,
	}
	return yaml.Marshal(res)
}
func (s *state) EventBytes() ([]byte, error) {
	res := map[string]interface{}{
		"events": s.events,
	}
	return yaml.Marshal(res)
}
func (s *state) Events() []event.Event {
	return s.events
}
func (s *state) Tasks() map[string]TaskState {
	return s.tasks
}
func (s *state) Services() []ServiceState {
	return s.services
}

func (s *state) StartProcess() {
	for {
		select {
		case _ = <-s.copyChanRequest:
			s.copyChan <- s.copy()
		case updateRequest := <-s.stateUpdateRequest:
			if updateRequest.AddRealtedTaskToEventReuqest != nil {
				s.logger.Debug("Updating state", "request", "AddRealtedTaskToEventReuqest")
				s.addRealtedTaskToEventReuqest(updateRequest.AddRealtedTaskToEventReuqest)
			}

			if updateRequest.ElectTasksRequest != nil {
				s.logger.Debug("Updating state", "request", "ElectTasksRequest")
				s.electTasksRequest(updateRequest.ElectTasksRequest)
			}

			if updateRequest.UpdateStateMetadataRequest != nil {
				s.logger.Debug("Updating state", "request", "UpdateStateMetadataRequest")
				s.updateStateMetadataRequest(updateRequest.UpdateStateMetadataRequest)
			}

			if updateRequest.UpdateTaskStateRequest != nil {
				s.logger.Debug("Updating state", "request", "UpdateTaskStateRequest")
				s.updateTaskStateRequest(updateRequest.UpdateTaskStateRequest)
			}
			s.wg.Done()
		}
	}
}

func (s *state) copy() state {
	destPtr := *s
	dest := destPtr
	dest.tasks = make(map[string]TaskState, len(s.tasks))
	for n, t := range s.tasks {
		dest.tasks[n] = t
	}
	dest.services = append([]ServiceState(nil), s.services...)
	dest.events = append([]event.Event(nil), s.events...)
	return dest
}

func (s *state) addRealtedTaskToEventReuqest(req *AddRealtedTaskToEventReuqest) {
	for i, e := range s.events {
		if e.Metadata.ID == req.EventID {
			event := event.Event{
				Metadata:     e.Metadata,
				Payload:      e.Payload,
				RelatedTasks: append(e.RelatedTasks, req.Task...),
			}
			// "update the event"
			s.events = append(append(s.events[0:i], event), s.events[i+1:]...)
			break
		}
	}
}

func (s *state) electTasksRequest(req *ElectTasksRequest) {
	for _, t := range req.Tasks {
		s.tasks[t.Name()] = TaskState{
			State: TaskStateElected,
			Task:  t,
		}
	}
	ev := &event.Event{
		Metadata: event.Metadata{
			CreatedAt: time.Now(),
			ID:        utils.GenerateID(),
			Name:      EventTaskElected,
		},
	}
	s.eventChan <- ev
}
func (s *state) updateTaskStateRequest(req *UpdateTaskStateRequest) {
	task := req.State.Task
	ev := &event.Event{
		Metadata: event.Metadata{
			CreatedAt: time.Now(),
			ID:        utils.GenerateID(),
			Task:      task.Name(),
		},
	}
	newstate := &TaskState{}
	if req.State.State == TaskStateInProgress {
		ev.Metadata.Name = EventTaskStarted
	} else if req.State.State == TaskStateFinished {
		ev.Metadata.Name = EventTaskFinished
	}
	t := s.tasks[task.Name()]
	if err := mergo.MergeWithOverwrite(newstate, t); err != nil {
		s.logger.Error("Failed to merge current task state into new task state", "err", err.Error())
	}
	if err := mergo.MergeWithOverwrite(newstate, req.State); err != nil {
		s.logger.Error("Failed to merge given task state into new task state", "err", err.Error())
	}
	s.tasks[task.Name()] = *newstate
	s.events = append(s.events, *ev)
	s.eventChan <- ev
}
func (s *state) updateStateMetadataRequest(req *UpdateStateMetadataRequest) {
	ev := &event.Event{
		Metadata: event.Metadata{
			CreatedAt: time.Now(),
			ID:        utils.GenerateID(),
		},
		RelatedTasks: []string{},
	}
	if req.State == EngineStateInProgress {
		s.state = EngineStateInProgress
		ev.Metadata.Name = EventEngineStarted
	} else if req.State == EngineStateFinished {
		s.state = EngineStateFinished
		ev.Metadata.Name = EventEngineFinished
	}
	s.events = append(s.events, *ev)
	s.eventChan <- ev
}
