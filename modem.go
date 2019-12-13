package core

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	v1 "github.com/open-integration/core/pkg/api/v1"
	"github.com/open-integration/core/pkg/logger"
	"github.com/open-integration/core/pkg/shell"
	"github.com/open-integration/core/pkg/utils"
	"github.com/xeipuuv/gojsonschema"
	"google.golang.org/grpc"
)

type (
	Modem interface {
		Init() error
		Call(spec TaskSpec, fd string) (string, error)
		Destory() error
		AddService(name string, port string, path string) error
	}

	modem struct {
		services            map[string]*service
		logger              logger.Logger
		serviceLogDirectory string
	}

	service struct {
		conn   *grpc.ClientConn
		client v1.ServiceClient
		ready  bool
		server struct {
			binPath string
			port    string
			pid     int
		}
		tasksSchemas map[string]string
	}

	ModemOptions struct {
		Logger logger.Logger
	}
)

func NewModem(opt *ModemOptions) Modem {
	m := &modem{
		logger:   opt.Logger,
		services: map[string]*service{},
	}
	return m
}

func (m *modem) Init() error {
	m.logger.Debug("Modem initializing started")
	for name, s := range m.services {
		log := m.logger.New("service", name)
		log.Debug("Starting service", "port", s.server.port)
		envs := []string{
			fmt.Sprintf("PORT=%s", s.server.port),
		}
		file, err := utils.CreateLogFile(m.serviceLogDirectory, fmt.Sprintf("%s-%d", name, time.Now().Unix()))
		if err != nil {
			log.Error("Failed to create log file", "error", err.Error())
			return err
		}
		log.Debug("Logging file created")
		logger := file
		detached := true
		pid, err := shell.Execute(s.server.binPath, []string{""}, envs, detached, "", "", logger)
		if err != nil {
			log.Error("Serivce startup failed", "error", err.Error())
			return err
		}
		s.server.pid = pid
		log.Debug("Server started")

		log.Debug("Dailing service")
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", s.server.port), grpc.WithInsecure())
		if err != nil {
			log.Error("Serivce dail failed", "error", err.Error())
			return err
		}
		log.Debug("Connection established")
		s.conn = conn
		client := v1.NewServiceClient(conn)
		s.client = client
		time.Sleep(2 * time.Second)
		log.Debug("Initializing service")
		resp, err := client.Init(context.Background(), &v1.InitRequest{})
		if err != nil {
			log.Error("Serivce init call failed", "error", err.Error())
			return err
		}
		s.tasksSchemas = resp.JsonSchemas

	}
	m.logger.Debug("Modem initializing finished")
	return nil
}
func (m *modem) Call(spec TaskSpec, fd string) (string, error) {
	m.logger.Debug("Calling", "service", spec.Service, "endpoint", spec.Endpoint)

	req := &v1.CallRequest{
		Endpoint: spec.Endpoint,
		Fd:       fd,
	}
	arguments := map[string]string{}
	for _, arg := range spec.Arguments {
		arguments[arg.Key] = arg.Value
	}
	err := m.isArgumentsValid(arguments, m.services[spec.Service].tasksSchemas[fmt.Sprintf("%s/%s", spec.Endpoint, "arguments.json")])
	if err != nil {
		return "", err
	}
	req.Arguments = arguments
	resp, err := m.services[spec.Service].client.Call(context.Background(), req)
	if err != nil {
		return "", err
	}
	if resp.Status == v1.Status_Error {
		return "", fmt.Errorf(resp.Error)
	}

	err = m.isResponsePayloadValid(resp.Payload, m.services[spec.Service].tasksSchemas[fmt.Sprintf("%s/%s", spec.Endpoint, "returns.json")])
	if err != nil {
		return "", err
	}

	return resp.Payload, nil
}

func (m *modem) Destory() error {
	m.logger.Debug("Stopping all services")
	for name, service := range m.services {
		if err := service.conn.Close(); err != nil {
			m.logger.Debug("Failed to close connection to service", "service", name)
		}
		process, err := os.FindProcess(service.server.pid)
		if err != nil {
			m.logger.Debug("Failed to find process of service", "service", name)
		}
		if err := process.Signal(os.Interrupt); err != nil {
			m.logger.Debug("Failed to send kill signal to service process", "service", name)
		}
		m.logger.Debug("Service stopped", "service", name)
	}
	return nil
}

func (m *modem) AddService(name string, port string, path string) error {
	s := &service{
		ready: false,
	}
	s.server.binPath = path
	s.server.port = port
	m.services[name] = s
	return nil
}

func (m *modem) isArgumentsValid(json map[string]string, schema string) error {
	if schema == "" {
		return nil // no schema given, no assertion required
	}
	schemaLoader := gojsonschema.NewStringLoader(schema)
	jsonLoader := gojsonschema.NewGoLoader(json)
	return m.isJsonValid(jsonLoader, schemaLoader)
}

func (m *modem) isResponsePayloadValid(json string, schema string) error {
	if schema == "" {
		return nil // no schema given, no assertion required
	}
	schemaLoader := gojsonschema.NewStringLoader(schema)
	jsonLoader := gojsonschema.NewStringLoader(json)
	return m.isJsonValid(jsonLoader, schemaLoader)
}

func (m *modem) isJsonValid(json gojsonschema.JSONLoader, schema gojsonschema.JSONLoader) error {
	result, err := gojsonschema.Validate(schema, json)
	if err != nil {
		return err
	}

	if !result.Valid() {
		message := strings.Builder{}
		for _, desc := range result.Errors() {
			message.WriteString(desc.String())
		}
		return fmt.Errorf(message.String())
	}
	return nil
}
