package engine

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/open-integration/oi/core/event"
	"github.com/open-integration/oi/core/event/reporter"
	"github.com/open-integration/oi/core/modem"
	"github.com/open-integration/oi/core/state"
	"github.com/open-integration/oi/core/task"

	"github.com/open-integration/oi/core/filedescriptor"
	"github.com/open-integration/oi/core/service/runner"
	"github.com/open-integration/oi/pkg/downloader"
	"github.com/open-integration/oi/pkg/graph"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/utils"
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
		eventChan          chan *event.Event
		stateUpdateRequest chan state.UpdateRequest
		taskLogsDirctory   string
		stateDir           string
		modem              modem.Modem
		statev1            state.State
		wg                 *sync.WaitGroup
		graphBuilder       graph.Builder
	}

	// Options to create new engine
	Options struct {
		Pipeline Pipeline
		// LogsDirectory path where to store logs
		LogsDirectory string
		Kubeconfig    *KubernetesOptions
		Logger        logger.Logger

		serviceDownloader downloader.Downloader
		modem             modem.Modem
	}

	// KubernetesOptions when running in/with kubernetes cluster
	KubernetesOptions struct {
		Path                string
		Context             string
		Namespace           string
		InCluster           bool
		Host                string
		B64Crt              string
		Token               string
		LogsVolumeClaimName string
		LogsVolumeName      string
	}
)

// New creates new Engine from options
func New(opt *Options) (Engine, error) {

	if opt.LogsDirectory == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("Failed to get current working directory: %w", err)
		}
		opt.LogsDirectory = wd
	}

	tasksLogDir := path.Join(opt.LogsDirectory, "logs", "tasks")
	if err := createDir(tasksLogDir); err != nil {
		return nil, fmt.Errorf("Failed to create task logs directory %s: %w", opt.LogsDirectory, err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("Failed to get user's home dir: %w", err)
	}

	servicesDir := path.Join(home, ".open-integration", "services")
	if err := createDir(servicesDir); err != nil {
		return nil, fmt.Errorf("Failed to create service cache directory %s: %w", servicesDir, err)
	}

	eventChannel := make(chan *event.Event, 10)

	stateUpdateChannel := make(chan state.UpdateRequest, 1)

	waitGroup := &sync.WaitGroup{}

	var log logger.Logger
	// create logger
	{
		if opt.Logger == nil {
			loggerOptions := &logger.Options{
				FilePath:    path.Join(opt.LogsDirectory, "logs", "log.log"),
				LogToStdOut: true,
			}
			log = logger.New(loggerOptions)
		} else {
			log = opt.Logger
		}
	}

	if opt.serviceDownloader == nil {
		opt.serviceDownloader = downloader.New(downloader.Options{
			Store:  servicesDir,
			Logger: log.New("module", "service-downloader"),
		})
	}

	servicesLogDir := path.Join(opt.LogsDirectory, "logs", "services")
	if err := createDir(servicesLogDir); err != nil {
		return nil, fmt.Errorf("Failed to create service logs directory %s: %w", servicesLogDir, err)
	}

	s := state.New(&state.Options{
		Logger:             log.New("module", "state-store"),
		EventChan:          eventChannel,
		CommandsChan:       make(chan string, 1),
		Name:               opt.Pipeline.Metadata.Name,
		StateUpdateRequest: stateUpdateChannel,
		WG:                 waitGroup,
	})
	go s.StartProcess()
	m, err := createModem(opt, log, servicesLogDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to create modemo: %w", err)
	}
	return &engine{
		statev1:            s,
		pipeline:           opt.Pipeline,
		wg:                 waitGroup,
		graphBuilder:       graph.New(),
		stateDir:           opt.LogsDirectory,
		taskLogsDirctory:   tasksLogDir,
		eventChan:          eventChannel,
		stateUpdateRequest: stateUpdateChannel,
		logger:             log.New("module", "engine"),
		modem:              m,
	}, nil
}

// Run starts the pipeline execution
func (e *engine) Run() error {
	var err error
	e.logger.Debug("Starting...", "pipeline", e.pipeline.Metadata.Name)
	err = e.modem.Init()
	if err != nil {
		return err
	}
	defer func() {
		e.logger.Debug("killing all services")
		if cerr := e.modem.Destroy(); cerr != nil {
			err = cerr
		}
	}()
	e.wg.Add(1)
	e.stateUpdateRequest <- state.UpdateRequest{
		Metadata: state.UpdateRequestMetadata{
			CreatedAt: utils.TimeNow(),
		},
		UpdateStateMetadataRequest: &state.UpdateStateMetadataRequest{
			State: state.EngineStateInProgress,
		},
	}
	go e.waitForFinish()
	e.handleStateEvents()
	e.printGraph()
	if cerr := e.printStateStore(); cerr != nil {
		return cerr
	}
	return err
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
		e.printGraph()
	}
}

// handleTaskEvents watch on dedicated event channel created for each task.
func (e *engine) handleTaskEvents(lgr logger.Logger, ch <-chan event.Event) {
	for {
		ev := <-ch
		lgr.Debug("Got event from task", "name", ev.Metadata.Name)
		go e.electNextTasks(ev)
	}
}

// electNextTasks - running all reactions on the event and sending request to elect matched tasks
func (e *engine) electNextTasks(ev event.Event) {
	log := e.logger.New("event", ev.Metadata.Name)
	log.Debug("Received event, electing next tasks")
	stateCpy, err := e.statev1.Copy()
	if err != nil {
		e.logger.Error("Failed to copy state")
		return
	}
	tasksCandidates := map[string]task.Task{}
	for _, reaction := range e.pipeline.Spec.Reactions {
		if reaction.Condition.Met(ev, stateCpy) {
			for _, t := range reaction.Reaction(ev, stateCpy) {
				tasksCandidates[t.Name()] = t
			}
		}
	}

	tasksToElect := []task.Task{}
	for _, t := range tasksCandidates {
		_, exist := stateCpy.Tasks()[t.Name()]
		if !exist {
			e.logger.Debug("Adding task to elected set", "task", t.Name())
			tasksToElect = append(tasksToElect, t)
		}
	}
	if len(tasksToElect) > 0 {
		e.logger.Debug("Electing tasks", "total", len(tasksToElect))
		ids := []string{}
		for _, t := range tasksToElect {
			ids = append(ids, t.Name())
		}
		e.wg.Add(1)
		e.stateUpdateRequest <- state.UpdateRequest{
			Metadata: state.UpdateRequestMetadata{
				CreatedAt: utils.TimeNow(),
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
func (e *engine) executeElectedTasks(ev event.Event) {
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
		log.Debug("Running task", "task", t.Name())
		go e.runTask(t, ev, log.New("task", t.Name()))
		wg.Done()
	}
	wg.Wait()
}

func (e *engine) runTask(t task.Task, ev event.Event, lgr logger.Logger) {
	fileName := fmt.Sprintf("%s.log", t.Name())
	fileDescriptor := path.Join(e.taskLogsDirctory, fileName)
	e.wg.Add(1)
	times := state.TaskTimes{
		Started: utils.TimeNow(),
	}
	e.stateUpdateRequest <- state.UpdateRequest{
		Metadata: state.UpdateRequestMetadata{
			CreatedAt: utils.TimeNow(),
		},
		UpdateTaskStateRequest: &state.UpdateTaskStateRequest{
			State: state.TaskState{
				State:  state.TaskStateInProgress,
				Task:   t,
				Times:  times,
				Logger: fileDescriptor,
			},
		},
	}

	_, err := utils.CreateLogFile(e.taskLogsDirctory, fileName)
	if err != nil {
		lgr.Error("Failed to create log file for task")
		return
	}

	fd, err := filedescriptor.New(fileDescriptor)
	if err != nil {
		lgr.Error("Failed to create filedescriptor")
		return
	}

	eventChan := make(chan event.Event)
	go e.handleTaskEvents(e.logger.New("module", "task-event-handler"), eventChan)
	payload, err := t.Run(context.Background(), task.RunOptions{
		FD: fd,
		EventReporter: reporter.New(reporter.Options{
			EventChan: eventChan,
			Name:      t.Name(),
		}),
		Modem: e.modem,
	})

	e.wg.Add(1)
	times.Finished = utils.TimeNow()
	e.stateUpdateRequest <- state.UpdateRequest{
		Metadata: state.UpdateRequestMetadata{
			CreatedAt: utils.TimeNow(),
		},
		UpdateTaskStateRequest: &state.UpdateTaskStateRequest{
			State: state.TaskState{
				State:  state.TaskStateFinished,
				Status: e.concludeStatus(err),
				Task:   t,
				Times:  times,
				Output: payload,
				Error:  err,
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
	e.stateUpdateRequest <- state.UpdateRequest{
		Metadata: state.UpdateRequestMetadata{
			CreatedAt: utils.TimeNow(),
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

func (e *engine) printGraph() {
	s, _ := e.statev1.Copy()
	g, err := e.graphBuilder.Build(s)
	if err != nil {
		e.logger.Error("Failed to build graph", "err", err.Error())
		return
	}
	if err := ioutil.WriteFile(path.Join(e.stateDir, "graph.dot"), g, os.ModePerm); err != nil {
		e.logger.Error("Failed to write graph to file", "err", err.Error())
	}

}

func (e *engine) concludeStatus(err error) string {
	status := state.TaskStatusSuccess
	if err != nil {
		e.logger.Error("Task exited with error", "err", err.Error())
		status = state.TaskStatusFailed
	}
	return status
}

func createDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func createModem(opt *Options, log logger.Logger, servicesLogDir string) (modem.Modem, error) {
	if opt.modem != nil {
		return opt.modem, nil
	}
	serviceModem := modem.New(&modem.Options{
		Logger: log.New("module", "modem"),
	})
	for _, s := range opt.Pipeline.Spec.Services {
		svcID := utils.GenerateID()
		if opt.Kubeconfig == nil {
			finalLocation := s.Path
			if s.Name != "" && s.Version != "" {
				location, err := opt.serviceDownloader.Download(s.Name, s.Version)
				if err != nil {
					return nil, fmt.Errorf("Failed to download serivce %s: %w", s.Name, err)
				}
				finalLocation = location
			}
			log.Debug("Adding service", "path", finalLocation)
			if err := serviceModem.AddService(s.As, runner.New(&runner.Options{
				Type:                 runner.Local,
				Logger:               log.New("service-service", s.Name),
				Name:                 s.Name,
				ID:                   svcID,
				Dailer:               &utils.GRPC{},
				PortGenerator:        utils.Port{},
				LocalLogFileCreator:  &utils.FileCreator{},
				LogsDirectory:        servicesLogDir,
				ServiceClientCreator: utils.Proto{},
				LocalCommandCreator:  &utils.Command{},
				LocalPathToBinary:    finalLocation,
			})); err != nil {
				return nil, fmt.Errorf("Failed to add service %s to modem: %w", s.Name, err)
			}
		} else {
			log.Debug("Adding service")
			runnerOpt := &runner.Options{
				Type:                      runner.Kubernetes,
				Logger:                    log.New("service-service", s.Name),
				Name:                      s.Name,
				ID:                        svcID,
				Version:                   s.Version,
				PortGenerator:             utils.Port{},
				KubernetesKubeConfigPath:  opt.Kubeconfig.Path,
				KubernetesContext:         opt.Kubeconfig.Context,
				KubernetesNamespace:       opt.Kubeconfig.Namespace,
				KubeconfigHost:            opt.Kubeconfig.Host,
				KubeconfigToken:           opt.Kubeconfig.Token,
				KubeconfigB64Crt:          opt.Kubeconfig.B64Crt,
				Kube:                      &utils.Kubernetes{},
				Dailer:                    &utils.GRPC{},
				ServiceClientCreator:      utils.Proto{},
				LogsDirectory:             opt.LogsDirectory,
				KubernetesVolumeClaimName: opt.Kubeconfig.LogsVolumeClaimName,
				KubernetesVolumeName:      opt.Kubeconfig.LogsVolumeName,
			}
			if opt.Kubeconfig.InCluster {
				runnerOpt.KubernetesGrpcDialViaPodIP = true
			}
			if err := serviceModem.AddService(s.As, runner.New(runnerOpt)); err != nil {
				return nil, fmt.Errorf("Failed to add service %s to modem: %w", s.Name, err)
			}
		}
	}
	return serviceModem, nil
}
