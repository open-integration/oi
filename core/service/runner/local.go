package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	v1 "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
	"google.golang.org/grpc"
)

type (
	localService struct {
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
		Create() *exec.Cmd
		AddCommand(cmd string)
		AddEnv(key string, value string)
		Bin(path string)
	}
)

func (_l *localService) Run() error {
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

func (_l *localService) Kill() error {
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

func (_l *localService) Call(context context.Context, req *v1.CallRequest) (*v1.CallResponse, error) {
	return _l.client.Call(context, req)
}

func (_l *localService) Schemas() map[string]string {
	return _l.tasksSchemas
}

func (_l *localService) generatePort() error {
	port, err := _l.portGenerator.GetAvailable()
	if err != nil {
		return err
	}
	_l.port = port
	return nil
}

func (_l *localService) generateLogFile() error {
	name := fmt.Sprintf("%s-%s.log", _l.name, _l.id)
	writer, err := _l.logFileCreator.Create(_l.logsFileDirectory, name)
	if err != nil {
		return err
	}
	_l.logWriter = writer
	return nil
}

func (_l *localService) createCommand() error {
	_l.cmdCreator.AddEnv("PORT", _l.port)
	_l.cmdCreator.Bin(_l.path)
	_l.command = _l.cmdCreator.Create()
	return nil
}

func (_l *localService) run() error {
	_l.command.Stdout = _l.logWriter
	_l.command.Stderr = _l.logWriter
	return _l.command.Start()
}

func (_l *localService) dail() error {
	url := fmt.Sprintf("localhost:%s", _l.port)
	conn, err := _l.dialer.Dial(url, grpc.WithInsecure())
	if err != nil {
		return err
	}
	_l.connection = conn
	_l.client = _l.serviceClientCreator.New(conn)
	_l.Logger.Debug("Connection established")
	return nil
}

func (_l *localService) init() error {
	_, err := _l.client.Init(context.Background(), &v1.InitRequest{})
	if err != nil {
		return err
	}
	return nil
}
