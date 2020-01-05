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

	var log logger.Logger
	{
		log = logger.New(loggerOptions)
	}

	serviceDownloader := downloader.New(downloader.Options{
		Store:  servicesDir,
		Logger: log.New("module", "service-downloader"),
	})

	servicesLogDir := path.Join(wd, "logs", "services")
	dieOnError(createDir(servicesLogDir))

	e.logger = log.New("module", "engine")

	// Init modem
	{
		e.modem = modem.New(&modem.ModemOptions{
			Logger: log.New("module", "modem"),
		})
		for _, s := range opt.Pipeline.Spec.Services {
			svcID := string(generateID())
			if opt.Kubeconfig == nil {
				location := s.Path
				if s.Name != "" && s.Version != "" {
					err := serviceDownloader.Download(s.Name, s.Version)
					dieOnError(err)
					location = path.Join(serviceDownloader.Store(), s.Name)
				}
				log.Debug("Adding service", "path", location)
				e.modem.AddService(svcID, s.As, runner.New(&runner.Options{
					Type:                 runner.LocalRunner,
					Logger:               log.New("service-runner", s.Name),
					Name:                 s.Name,
					ID:                   svcID,
					Dailer:               &utils.GRPC{},
					PortGenerator:        utils.Port{},
					LocalLogFileCreator:  &utils.FileCreator{},
					LocalLogsDirectory:   servicesLogDir,
					ServiceClientCreator: utils.Proto{},
					LocalCommandCreator:  utils.Command{},
					LocalPathToBinary:    location,
				}))
			} else {
				log.Debug("Adding service")
				e.modem.AddService(svcID, s.As, runner.New(&runner.Options{
					Type:                     runner.KubernetesRunner,
					Logger:                   log.New("service-runner", s.Name),
					Name:                     s.Name,
					ID:                       svcID,
					Version:                  s.Version,
					PortGenerator:            utils.Port{},
					KubernetesKubeConfigPath: opt.Kubeconfig.Path,
					KubernetesContext:        opt.Kubeconfig.Context,
					KubernetesNamespace:      opt.Kubeconfig.Namespace,
					Kube:                     &utils.Kubernetes{},
					Dailer:                   &utils.GRPC{},
					ServiceClientCreator:     utils.Proto{},
				}))
			}
		}
	}

	e.state = NewState(&StateOptions{
		Logger:           e.logger.New("sub-module", "state-store"),
		EventCn:          eventCn,
		StateFile:        "./logs/state.yaml",
		EventHistoryFile: "./logs/history.yaml",
	})
	return e
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
