package core

import (
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/imdario/mergo"
	"github.com/open-integration/core/internal/commands"
	"github.com/open-integration/core/pkg/logger"
	"gopkg.in/yaml.v2"
)

const (
	// EngineStateInProgress pipeline in execution progress
	EngineStateInProgress string = "in-progress"
	// EngineStateFinished pipeline is finished execution
	EngineStateFinished string = "finished"

	// TaskStateInProgress task is in progress
	TaskStateInProgress string = EngineStateInProgress

	// TaskStateFinished task is finished
	TaskStateFinished string = EngineStateFinished

	// TaskStatusSuccess set one the task status in case task was finished successfully
	TaskStatusSuccess = "Success"
	// TaskStatusFailed set one the task status in case task was finished and failed
	TaskStatusFailed = "failed"
)

type (

	// StateOptions to creates state
	StateOptions struct {
		EventCn          chan *Event
		Logger           logger.Logger
		StateFile        string
		EventHistoryFile string
	}

	historyElem struct {
		Event Event  `yaml:"event"`
		State *State `yaml:"state"`
	}

	// State is the overall state of the pipeline execution
	State struct {
		eventChan        chan *Event
		History          []historyElem `yaml:"history"`
		logger           logger.Logger
		Metadata         StateMetadata    `yaml:"metadata"`
		Tasks            map[ID]TaskState `yaml:"tasks"`
		stateFile        string
		eventHistoryFile string
	}

	// ID is Uniqe ID
	ID string

	// StateMetadata metadata of the statee
	StateMetadata struct {
		// State is in-progress / finished
		State string `yaml:"state"`
	}

	// TaskState state of singal task
	TaskState struct {
		State   string `yaml:"state"`
		Status  string `yaml:"status"`
		Task    string `yaml:"task"`
		Output  string `yaml:"output"`
		Error   string `yaml:"error"`
		EventID ID     `yaml:"event-id"`
	}
)

// NewState creates new state store
func NewState(opt *StateOptions) *State {
	return &State{
		eventChan:        opt.EventCn,
		logger:           opt.Logger,
		stateFile:        opt.StateFile,
		eventHistoryFile: opt.EventHistoryFile,
	}
}

func (s *State) send(cmd commands.Command, wg *sync.WaitGroup, update func() *State) error {
	s.logger.Debug("Received command", "command", cmd.String())
	var ev Event
	switch cmd {
	case commands.StartEngine:
		ev = Event{
			Metadata: EventMetadata{
				Name:      EventEngineStarted,
				ID:        generateID(),
				CreatedAt: time.Now(),
			},
		}
		s.eventChan <- &ev
	case commands.FinishEngine:
		ev = Event{
			Metadata: EventMetadata{
				Name:      EventEngineFinished,
				ID:        generateID(),
				CreatedAt: time.Now(),
			},
		}
		s.eventChan <- &ev
	case commands.StartTask:
		ev = Event{
			Metadata: EventMetadata{
				Name:      EventTaskStarted,
				ID:        generateID(),
				CreatedAt: time.Now(),
			},
		}
		s.eventChan <- &ev
	case commands.FinishTask:
		ev = Event{
			Metadata: EventMetadata{
				Name:      EventTaskFinished,
				ID:        generateID(),
				CreatedAt: time.Now(),
			},
		}
		s.eventChan <- &ev
	}

	var data *State = nil
	if update != nil {
		data = update()
	}
	s.History = append(s.History, historyElem{
		Event: ev,
		State: data,
	})

	if data != nil {
		if err := mergo.Merge(s, data, mergo.WithOverride); err != nil {
			s.logger.Error("Failed to update state", "msg", err.Error())
			return err
		}
	}
	wg.Done()
	return s.printStore()
}

func (s *State) printStore() error {
	store := map[string]interface{}{
		"metadata": s.Metadata,
		"tasks":    s.Tasks,
	}
	storeJSON, err := yaml.Marshal(store)
	if err != nil {
		s.logger.Error("Failed to marshal state")
		return err
	}

	events := map[string]interface{}{
		"events": s.History,
	}
	eventsJSON, err := yaml.Marshal(events)
	if err != nil {
		s.logger.Error("Failed to marshal state")
		return err
	}
	if s.stateFile != "" {
		err := ioutil.WriteFile(s.stateFile, storeJSON, os.ModePerm)
		if err != nil {
			s.logger.Error("Failed to store state to file")
			return err
		}

	}

	if s.eventHistoryFile != "" {
		err = ioutil.WriteFile(s.eventHistoryFile, eventsJSON, os.ModePerm)
		if err != nil {
			s.logger.Error("Failed to store state to file")
			return err
		}
	}
	return nil
}
