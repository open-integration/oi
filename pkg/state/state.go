package state

import (
	"errors"
	"sync"
	"time"

	"github.com/imdario/mergo"
	"github.com/open-integration/core/pkg/logger"
	"github.com/open-integration/core/pkg/utils"
	"gopkg.in/yaml.v2"
)

type (
	// State holds all the data of the pipeline execution flow
	State interface {
		Copy() (State, error)
		Tasks() map[string]TaskState
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
		events             []Event
		eventChan          chan *Event
		commandsChan       chan string
		logger             logger.Logger
		initialized        bool
		copyChan           chan state
		copyChanRequest    chan bool
		stateUpdateRequest chan StateUpdateRequest
		wg                 *sync.WaitGroup
	}

	// Options to pass to the state
	Options struct {
		Name string
		// EventChan to write new event to the channel once a chance was applied
		EventChan chan *Event
		// CommandsChan to receive commands to create new change channel
		CommandsChan chan string
		// StateUpdateRequest to receive updated on the state in realtime
		StateUpdateRequest chan StateUpdateRequest
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
			"status": s.state,
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
				continue
			}

			if updateRequest.ElectTasksRequest != nil {
				s.logger.Debug("Updating state", "request", "ElectTasksRequest")
				s.electTasksRequest(updateRequest.ElectTasksRequest)
				continue
			}

			if updateRequest.UpdateStateMetadataRequest != nil {
				s.logger.Debug("Updating state", "request", "UpdateStateMetadataRequest")
				s.updateStateMetadataRequest(updateRequest.UpdateStateMetadataRequest)
				continue
			}

			if updateRequest.UpdateTaskStateRequest != nil {
				s.logger.Debug("Updating state", "request", "UpdateTaskStateRequest")
				s.updateTaskStateRequest(updateRequest.UpdateTaskStateRequest)
				continue
			}
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
	return dest
}

func (s *state) addRealtedTaskToEventReuqest(req *AddRealtedTaskToEventReuqest) {
	for i, e := range s.events {
		if e.Metadata.ID == req.EventID {
			event := Event{
				Metadata:     e.Metadata,
				Payload:      e.Payload,
				RelatedTasks: append(e.RelatedTasks, req.Task),
			}
			s.events = append(append(s.events[0:i], event), s.events[i+1:]...)
			break
		}
	}
}

func (s *state) electTasksRequest(req *ElectTasksRequest) {
	for _, id := range req.Tasks {
		s.tasks[id] = TaskState{
			State: TaskStateInProgress,
		}
	}
}
func (s *state) updateTaskStateRequest(req *UpdateTaskStateRequest) {
	ev := &Event{
		Metadata: EventMetadata{
			CreatedAt: time.Now(),
			ID:        utils.GenerateID(),
		},
	}
	id := req.State.Task
	if req.State.State == TaskStateInProgress {
		ev.Metadata.Name = EventTaskStarted
	} else if req.State.State == TaskStateFinished {
		ev.Metadata.Name = EventTaskFinished
	}
	newstate := &TaskState{}
	mergo.MergeWithOverwrite(newstate, s.tasks[id])
	mergo.MergeWithOverwrite(newstate, req.State)
	s.tasks[id] = *newstate
	s.events = append(s.events, *ev)
	s.eventChan <- ev
	s.wg.Done()
}
func (s *state) updateStateMetadataRequest(req *UpdateStateMetadataRequest) {
	ev := &Event{
		Metadata: EventMetadata{
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
	s.wg.Done()
}
