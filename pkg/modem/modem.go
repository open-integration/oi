package modem

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
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
		Call(service string, endpoint string, arguments map[string]interface{}, fd string) (string, error)
		Destory() error
		AddService(id string, name string, port string, path string) error
	}

	modem struct {
		services            map[string]*service
		logger              logger.Logger
		serviceLogDirectory string
		wg                  sync.WaitGroup
	}

	service struct {
		conn   *grpc.ClientConn
		client v1.ServiceClient
		ready  bool
		id     string
		server struct {
			binPath string
			port    string
			pid     int
		}
		err          error
		tasksSchemas map[string]string
	}

	ModemOptions struct {
		Logger              logger.Logger
		ServiceLogDirectory string
	}
)

func New(opt *ModemOptions) Modem {
	m := &modem{
		logger:              opt.Logger,
		services:            make(map[string]*service),
		serviceLogDirectory: opt.ServiceLogDirectory,
	}
	return m
}

func (m *modem) Init() error {
	m.logger.Debug("Modem initializing started")
	for name, s := range m.services {
		m.wg.Add(1)
		go m.initService(name, s, m.logger.New("service", name))
	}
	m.wg.Wait()

	for name, s := range m.services {
		if s.err != nil {
			m.logger.Crit("Service failed to start", "service", name, "err", s.err.Error())
		}
	}
	m.logger.Debug("Modem initializing finished")
	return nil
}

func (m *modem) Call(service string, endpoint string, arguments map[string]interface{}, fd string) (string, error) {
	log := m.logger.New("service", service, "endpoint", endpoint)
	log.Debug("Call service request")

	req := &v1.CallRequest{
		Endpoint: endpoint,
		Fd:       fd,
	}
	argsJSON, err := json.Marshal(arguments)
	if err != nil {
		return "", err
	}
	log.Debug("Validating arguments")
	err = m.isArgumentsValid(argsJSON, m.services[service].tasksSchemas[fmt.Sprintf("%s/%s", endpoint, "arguments.json")])
	if err != nil {
		return "", err
	}
	log.Debug("Arguments are valid")
	req.Arguments = string(argsJSON)
	resp, err := m.services[service].client.Call(context.Background(), req)
	if err != nil {
		log.Debug("Call return with error", "err", err.Error())
		return "", err
	}
	if resp.Status == v1.Status_Error {
		log.Debug("Call return with error", "err", resp.Error)
		return resp.Payload, fmt.Errorf(resp.Error)
	}

	log.Debug("Call ended", "response", resp.Status)
	err = m.isResponsePayloadValid(resp.Payload, m.services[service].tasksSchemas[fmt.Sprintf("%s/%s", endpoint, "returns.json")])
	if err != nil {
		return resp.Payload, err
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

func (m *modem) AddService(id string, name string, port string, path string) error {
	s := &service{
		ready: false,
		id:    id,
	}
	s.server.binPath = path
	s.server.port = port
	m.services[name] = s
	return nil
}

func (m *modem) initService(name string, svc *service, log logger.Logger) {
	log.Debug("Starting service", "port", svc.server.port)
	envs := []string{
		fmt.Sprintf("PORT=%s", svc.server.port),
	}
	logFile := fmt.Sprintf("%s-%s.log", name, svc.id)
	file, err := utils.CreateLogFile(m.serviceLogDirectory, logFile)
	if err != nil {
		log.Error("Failed to create log file", "error", err.Error())
		svc.err = err
		m.wg.Done()
		return
	}
	log.Debug("Logging file created", "file", name)
	logger := file
	detached := true
	pid, err := shell.Execute(svc.server.binPath, []string{""}, envs, detached, "", "", logger)
	if err != nil {
		log.Error("Serivce startup failed", "error", err.Error())
		svc.err = err
		m.wg.Done()
		return
	}
	svc.server.pid = pid
	log.Debug("Server started")

	log.Debug("Dailing service")
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", svc.server.port), grpc.WithInsecure())
	if err != nil {
		log.Error("Serivce dail failed", "error", err.Error())
		svc.err = err
		m.wg.Done()
		return
	}
	log.Debug("Connection established")
	svc.conn = conn
	client := v1.NewServiceClient(conn)
	svc.client = client
	time.Sleep(2 * time.Second)
	log.Debug("Initializing service")
	resp, err := client.Init(context.Background(), &v1.InitRequest{})
	if err != nil {
		log.Error("Serivce init call failed", "error", err.Error())
		svc.err = err
		m.wg.Done()
		return
	}
	svc.tasksSchemas = resp.JsonSchemas
	m.wg.Done()
	return
}

func (m *modem) isArgumentsValid(json []byte, schema string) error {
	if schema == "" {
		return nil // no schema given, no assertion required
	}
	schemaLoader := gojsonschema.NewStringLoader(schema)
	jsonLoader := gojsonschema.NewBytesLoader(json)
	return m.isJSONValid(jsonLoader, schemaLoader)
}

func (m *modem) isResponsePayloadValid(json string, schema string) error {
	if schema == "" {
		return nil // no schema given, no assertion required
	}
	schemaLoader := gojsonschema.NewStringLoader(schema)
	jsonLoader := gojsonschema.NewStringLoader(json)
	return m.isJSONValid(jsonLoader, schemaLoader)
}

func (m *modem) isJSONValid(json gojsonschema.JSONLoader, schema gojsonschema.JSONLoader) error {
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
