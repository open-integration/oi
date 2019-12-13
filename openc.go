package openc

import (
	"fmt"
	"os"
	"path"

	"github.com/open-integration/core/pkg/downloader"
	"github.com/open-integration/core/pkg/logger"
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
		FilePath: path.Join(wd, "logs", "openc.log"),
	}

	home, err := os.UserHomeDir()
	dieOnError(err)

	servicesDir := path.Join(home, ".openc", "services")
	dieOnError(createDir(servicesDir))

	serviceDownloader := downloader.New(downloader.Options{
		Store:  servicesDir,
		Logger: logger.New(loggerOptions).New("module", "service-downloader"),
	})

	servicesLogDir := path.Join(wd, "logs", "tasks")
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

func newModem(pipeline *Pipeline, servicesLogDir string, downloader downloader.Downloader, logger logger.Logger) Modem {

	m := &modem{
		logger:              logger,
		services:            map[string]*service{},
		serviceLogDirectory: servicesLogDir,
	}
	for _, p := range pipeline.Spec.Services {
		err := downloader.Download(p.Name, p.Version)
		dieOnError(err)
		port, err := utils.GetAvailablePort()
		dieOnError(err)
		logger.Debug("Adding service", "path", path.Join(downloader.Store(), p.Name))
		m.AddService(p.As, port, path.Join(downloader.Store(), p.Name))
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
