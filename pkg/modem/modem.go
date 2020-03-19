package modem

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	v1 "github.com/open-integration/core/pkg/api/v1"
	"github.com/open-integration/core/pkg/logger"
	"github.com/open-integration/core/pkg/runner"
	"github.com/xeipuuv/gojsonschema"
)

type (
	Modem interface {
		Init() error
		Call(service string, endpoint string, arguments map[string]interface{}, fd string) (string, error)
		Destroy() error
		AddService(id string, name string, runner runner.Runner) error
	}

	modem struct {
		services map[string]*service
		logger   logger.Logger
		wg       *sync.WaitGroup
	}

	service struct {
		ready        bool
		id           string
		runner       runner.Runner
		err          error
		tasksSchemas map[string]string
	}

	ModemOptions struct {
		Logger logger.Logger
	}
)

func New(opt *ModemOptions) Modem {
	m := &modem{
		logger:   opt.Logger,
		services: make(map[string]*service),
		wg:       &sync.WaitGroup{},
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

	err := strings.Builder{}
	for name, s := range m.services {
		if s.err != nil {
			m.logger.Error("Service failed to start", "service", name, "err", s.err.Error())
			err.WriteString(fmt.Sprintf("Serive: %s - Error: %s\n", name, s.err.Error()))
		}
	}
	if err.Len() > 0 {
		m.logger.Error("Modem initializing finished with error")
		return errors.New(err.String())
	}
	m.logger.Debug("Modem initializing finished")
	return nil
}

func (m *modem) Call(service string, endpoint string, arguments map[string]interface{}, fd string) (string, error) {
	log := m.logger.New("service", service, "endpoint", endpoint)

	if _, ok := m.services[service]; !ok {
		return "", fmt.Errorf("Service %s not found", service)
	}

	req := &v1.CallRequest{
		Endpoint: endpoint,
		Fd:       fd,
	}
	schema, ok := m.services[service].tasksSchemas[fmt.Sprintf("%s/%s", endpoint, "arguments.json")]
	if !ok {
		req.Arguments = ""
	} else {
		r, err := json.Marshal(arguments)
		if err != nil {
			return "", err
		}
		err = m.isArgumentsValid(r, schema)
		if err != nil {
			return "", err
		}
		req.Arguments = string(r)
	}
	resp, err := m.services[service].runner.Call(context.Background(), req)
	if err != nil {
		log.Debug("Call return with error", "err", err.Error())
		return "", err
	}
	if resp.Status == v1.Status_Error {
		log.Debug("Call return with error", "err", resp.Error)
		return resp.Payload, fmt.Errorf(resp.Error)
	}

	err = m.isResponsePayloadValid(resp.Payload, m.services[service].tasksSchemas[fmt.Sprintf("%s/%s", endpoint, "returns.json")])
	if err != nil {
		return resp.Payload, err
	}

	return resp.Payload, nil
}

func (m *modem) Destroy() error {
	for name, service := range m.services {
		err := service.runner.Kill()
		if err != nil {
			m.logger.Debug("Failed to kill service", "service", name)
		}
		m.logger.Debug("Service stopped", "service", name)
	}
	return nil
}

func (m *modem) AddService(id string, name string, runner runner.Runner) error {
	s := &service{
		ready: false,
		id:    id,
	}
	s.runner = runner
	m.services[name] = s
	return nil
}

func (m *modem) initService(name string, svc *service, log logger.Logger) {
	defer m.wg.Done()
	if err := svc.runner.Run(); err != nil {
		log.Error("Serivce startup failed", "error", err.Error())
		svc.err = err
		return
	}
	svc.ready = true
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
