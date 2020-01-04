package core

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/open-integration/core/pkg/downloader"
	"github.com/open-integration/core/pkg/logger"
	"github.com/open-integration/core/pkg/modem"
	"github.com/open-integration/core/pkg/runner"
	"github.com/open-integration/core/pkg/utils"
)

// NewEngine create new engine
func NewEngine(opt *EngineOptions) Engine {
	wd, err := os.Getwd()
	dieOnError(err)

	tasksLogDir := path.Join(wd, "logs", "tasks")
	dieOnError(createDir(tasksLogDir))

	eventCn := make(chan *Event, 1)
	e := &engine{
		pipeline:         opt.Pipeline,
		eventChan:        eventCn,
		taskLogsDirctory: tasksLogDir,
	}

	var loggerOptions *logger.Options

	loggerOptions = &logger.Options{
		FilePath:    path.Join(wd, "logs", "log"),
		LogToStdOut: true,
	}

	home, err := os.UserHomeDir()
	dieOnError(err)

	servicesDir := path.Join(home, ".open-integration", "services")
	dieOnError(createDir(servicesDir))

	serviceDownloader := downloader.New(downloader.Options{
		Store:  servicesDir,
		Logger: logger.New(loggerOptions).New("module", "service-downloader"),
	})

	servicesLogDir := path.Join(wd, "logs", "services")
	dieOnError(createDir(servicesLogDir))

	if opt.Logger == nil {
		e.logger = logger.New(loggerOptions).New("module", "engine")
		e.modem = newModem(&opt.Pipeline, servicesLogDir, serviceDownloader, logger.New(loggerOptions).New("module", "modem"))
	}
	e.state = NewState(&StateOptions{
		Logger:           e.logger.New("sub-module", "state-store"),
		EventCn:          eventCn,
		StateFile:        "./logs/state.yaml",
		EventHistoryFile: "./logs/history.yaml",
	})
	return e
}

func newModem(pipeline *Pipeline, servicesLogDir string, downloader downloader.Downloader, log logger.Logger) modem.Modem {
	m := modem.New(&modem.ModemOptions{
		Logger: log,
	})
	for _, s := range pipeline.Spec.Services {
		location := s.Path
		if s.Name != "" && s.Version != "" {
			err := downloader.Download(s.Name, s.Version)
			dieOnError(err)
			location = path.Join(downloader.Store(), s.Name)
		}
		log.Debug("Adding service", "path", location)
		svcID := string(generateID())
		m.AddService(svcID, s.As, runner.New(&runner.Options{
			Type: runner.LocalRunner,
			// TODO: remove module=modem
			Logger:                    log.New("service-runner", s.Name),
			Name:                      s.Name,
			ID:                        svcID,
			LocalDailer:               &utils.GRPC{},
			LocalLogFileCreator:       &utils.FileCreator{},
			LocalLogsDirectory:        servicesLogDir,
			LocalServiceClientCreator: utils.Proto{},
			LocalPortGenerator:        utils.Port{},
			LocalCommandCreator:       utils.Command{},
			LocalPathToBinary:         location,
		}))
	}
	return m
}

func dieOnError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func createDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func buildServiceCmd(svc Service, port string, location string) *exec.Cmd {
	cmd := exec.Command(location)
	envs := []string{
		fmt.Sprintf("PORT=%s", port),
	}
	cmd.Env = envs
	cmd.Dir = ""
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}
