package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	v1 "github.com/open-integration/core/pkg/api/v1"
	"github.com/open-integration/core/pkg/logger"
	"google.golang.org/grpc"
)

type (
	localRunner struct {
		Logger               logger.Logger
		command              *exec.Cmd
		name                 string
		id                   string
		path                 string
		port                 string
		logFileCreator       logFileCreator
		logsFileDirectory    string
		logWriter            io.Writer
		dialer               dialer
		connection           *grpc.ClientConn
		client               v1.ServiceClient
		serviceClientCreator serviceClientCreator
		tasksSchemas         map[string]string
		portGenerator        portGenerator
		cmdCreator           cmdCreator
	}

	cmdCreator interface {
		Create(port string, path string) *exec.Cmd
	}
)

func (_l *localRunner) Run() error {
	_l.Logger.Debug("Initializing service")
	if err := _l.generatePort(); err != nil {
		return err
	}

	if err := _l.generateLogFile(); err != nil {
		return err
	}

	if err := _l.createCommand(); err != nil {
		return err
	}

	if err := _l.run(); err != nil {
		return err
	}

	if err := _l.dail(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	if err := _l.init(); err != nil {
		return err
	}
	return nil
}

func (_l *localRunner) Kill() error {
	_l.Logger.Debug("Killing service")

	if err := _l.connection.Close(); err != nil {
		return err
	}

	process, err := os.FindProcess(_l.command.Process.Pid)
	if err != nil {
		return err
	}
	return process.Signal(os.Interrupt)
}

func (_l *localRunner) Call(context context.Context, req *v1.CallRequest) (*v1.CallResponse, error) {
	return _l.client.Call(context, req)
}

func (_l *localRunner) generatePort() error {
	_l.Logger.Debug("Generating available port on machine")
	port, err := _l.portGenerator.GetAvailable()
	if err != nil {
		return err
	}
	_l.port = port
	_l.Logger.Debug("Port generated", "port", port)
	return nil
}

func (_l *localRunner) generateLogFile() error {
	name := fmt.Sprintf("%s-%s.log", _l.name, _l.id)
	_l.Logger.Debug("Creating log file", "file", name)
	writer, err := _l.logFileCreator.Create(_l.logsFileDirectory, name)
	if err != nil {
		return err
	}
	_l.logWriter = writer
	_l.Logger.Debug("Log file created", "file", name)
	return nil
}

func (_l *localRunner) createCommand() error {
	_l.command = _l.cmdCreator.Create(_l.port, _l.path)
	return nil
}

func (_l *localRunner) run() error {
	_l.Logger.Debug("Starting service binary")
	_l.command.Stdout = _l.logWriter
	_l.command.Stderr = _l.logWriter
	return _l.command.Start()
}

func (_l *localRunner) dail() error {
	url := fmt.Sprintf("localhost:%s", _l.port)
	_l.Logger.Debug("Dial to service", "URL", url)
	conn, err := _l.dialer.Dial(url, grpc.WithInsecure())
	if err != nil {
		return err
	}
	_l.connection = conn
	_l.client = _l.serviceClientCreator.New(conn)
	_l.Logger.Debug("Connection established")
	return nil
}

func (_l *localRunner) init() error {
	_l.Logger.Debug("Calling service init endpoint one time")
	resp, err := _l.client.Init(context.Background(), &v1.InitRequest{})
	if err != nil {
		return err
	}
	_l.tasksSchemas = resp.JsonSchemas
	return nil
}
