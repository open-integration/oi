package state

import (
	"errors"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"github.com/open-integration/core/pkg/logger"
	"gopkg.in/yaml.v2"
)

var channel = 0

type (
	// State holds all the data of the pipeline execution flow
	State interface {
		Copy() (State, error)
		Tasks() []TaskState
		Services() []ServiceState
		Bytes() ([]byte, error)
		Send(cmd string, wg *sync.WaitGroup) chan Change
	}

	state struct {
		name                     string
		state                    string
		changes                  []Change
		services                 []ServiceState
		tasks                    []TaskState
		eventChan                chan *Event
		changesChan              chan Change
		commandsChan             chan string
		logger                   logger.Logger
		commandChangePairChannel chan commandChangePair
		initialized              bool
		copyChan                 chan state
		copyChanRequest          chan bool
	}

	// Options to pass to the state
	Options struct {
		Name string
		// EventChan to write new event to the channel once a chance was applied
		EventChan chan *Event
		// CommandsChan to receive commands to create new change channel
		CommandsChan chan string
		Logger       logger.Logger
	}

	commandChangePair struct {
		cmd    string
		change Change
	}
)

// New builds nnew State from options
func New(opt *Options) State {
	return &state{
		eventChan:                opt.EventChan,
		commandsChan:             opt.CommandsChan,
		logger:                   opt.Logger,
		commandChangePairChannel: make(chan commandChangePair, 1),
		initialized:              false,
		copyChan:                 make(chan state, 1),
		copyChanRequest:          make(chan bool, 1),
	}
}

// Copy
func (s *state) Copy() (State, error) {
	s.copyChanRequest <- true
	for {
		select {
		case cpy, _ := <-s.copyChan:
			return &cpy, nil
		case <-time.After(10 * time.Second):
			msg := "Failed to copy state after 10 seconds"
			return nil, errors.New(msg)
		}
	}
}
func (s *state) Bytes() ([]byte, error) {
	res := map[string]interface{}{
		"metadata": map[string]interface{}{
			"status": s.state,
		},
		"tasks":   s.tasks,
		"service": s.services,
		"changes": s.changes,
	}
	return yaml.Marshal(res)
}
func (s *state) Tasks() []TaskState {
	return s.tasks
}
func (s *state) Services() []ServiceState {
	return s.services
}

func (s *state) Send(cmd string, wg *sync.WaitGroup) chan Change {
	log := s.logger.New("command", cmd)
	log.Debug("Received command, creating disposable channel")
	if !s.initialized {
		log.Debug("Received command first time, running state initialization process")
		go s.init()
	}
	ch := make(chan Change)
	go func() {
		log.Debug("Wating to recieve change on created channel")
		change := <-ch
		log.Debug("Received change, sending to root change channel")
		s.commandChangePairChannel <- commandChangePair{
			cmd:    cmd,
			change: change,
		}
		log.Debug("Closeing disposable channel")
		close(ch)
		wg.Done()
	}()
	return ch
}

func (s *state) init() {
	s.initialized = true
	for {
		select {
		case pair := <-s.commandChangePairChannel:
			s.logger.Debug("Received command-change")
			ev := &Event{
				Metadata: EventMetadata{
					CreatedAt: time.Now(),
				},
			}
			s.changes = append(s.changes, pair.change)
			switch pair.cmd {
			case CmdStartEngine:
				ev.Metadata.Name = EventEngineStarted
				s.state = pair.change.PipelineChange.States
			case CmdFinishEngine:
				ev.Metadata.Name = EventEngineFinished
				s.state = pair.change.PipelineChange.States
			case CmdStartTask:
				ev.Metadata.Name = EventTaskStarted
				s.tasks = append(s.tasks, s.taskStateFromCommandChange(pair))
			case CmdFinishTask:
				newTaskState := s.taskStateFromCommandChange(pair)
				newSet := []TaskState{}
				for i, t := range s.tasks {
					if t.ID == pair.change.TaskChanges.ID {
						newSet = append(append(s.tasks[0:i], newTaskState), s.tasks[i+1:]...)
						break
					}
				}
				s.logger.Debug("Task set updated", "len", len(s.tasks))
				s.tasks = newSet
				ev.Metadata.Name = EventTaskFinished
			}

			s.logger.Debug("Sending event", "event", ev)
			s.eventChan <- ev
		case _ = <-s.copyChanRequest:
			dest := &state{}
			copier.Copy(dest, s)
			s.copyChan <- *dest
			continue
		}
	}
}

func (s *state) taskStateFromCommandChange(pair commandChangePair) TaskState {
	return TaskState{
		ID:        pair.change.TaskChanges.ID,
		Arguments: pair.change.TaskChanges.Arguments,
		Output:    pair.change.TaskChanges.Output,
		Error:     pair.change.TaskChanges.Error,
		EventID:   pair.change.TaskChanges.EventID,
		Logger:    pair.change.TaskChanges.LoggerID,
		State:     pair.change.TaskChanges.State,
		Status:    pair.change.TaskChanges.Status,
		Task:      pair.change.TaskChanges.Name,
	}
}
