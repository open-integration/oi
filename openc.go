package core

import (
	"fmt"
	"os"
	"path"

	"github.com/open-integration/core/pkg/downloader"
	"github.com/open-integration/core/pkg/logger"
	"github.com/open-integration/core/pkg/modem"
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

func newModem(pipeline *Pipeline, servicesLogDir string, downloader downloader.Downloader, logger logger.Logger) modem.Modem {
	m := modem.New(&modem.ModemOptions{
		Logger:              logger,
		ServiceLogDirectory: servicesLogDir,
	})
	for _, p := range pipeline.Spec.Services {
		location := p.Path
		if p.Name != "" && p.Version != "" {
			err := downloader.Download(p.Name, p.Version)
			dieOnError(err)
			location = path.Join(downloader.Store(), p.Name)
		}
		port, err := utils.GetAvailablePort()
		dieOnError(err)
		logger.Debug("Adding service", "path", location)
		m.AddService(string(generateID()), p.As, port, location)
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
