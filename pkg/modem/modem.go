package modem

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	v1 "github.com/open-integration/core/pkg/api/v1"
	"github.com/open-integration/core/pkg/logger"
	"github.com/open-integration/core/pkg/service"
	"github.com/xeipuuv/gojsonschema"
)

type (
	// Modem holds all the connection the external services
	Modem interface {
		ServiceCaller
		Init() error
		Destroy() error
		AddService(name string, runner service.Service) error
	}

	modem struct {
		services map[string]service.Service
		logger   logger.Logger
	}

	// Options to build new Modem
	Options struct {
		Logger logger.Logger
	}
)

// New build Modem from options
func New(opt *Options) Modem {
	m := &modem{
		logger:   opt.Logger,
		services: map[string]service.Service{},
	}
	return m
}

// Init starts the modem and all the services
func (m *modem) Init() error {
	m.logger.Debug("Modem initiation started")
	for name, s := range m.services {
		if err := m.initService(name, s, m.logger.New("service", name)); err != nil {
			return fmt.Errorf("Failed to initiate service: %w", err)
		}
	}
	m.logger.Debug("Modem initiation finished")
	return nil
}

// Call calls a service with input and returns output
func (m *modem) Call(ctx context.Context, service string, endpoint string, arguments map[string]interface{}, fd string) ([]byte, error) {
	log := m.logger.New("service", service, "endpoint", endpoint)
	req := &v1.CallRequest{
		Endpoint: endpoint,
		Fd:       fd,
	}
	svc, ok := m.services[service]
	if !ok {
		return nil, fmt.Errorf("Service %s not found", service)
	}
	schemas := svc.Schemas()
	schema, ok := schemas[fmt.Sprintf("endpoints/%s/%s", endpoint, "arguments.json")]
	if !ok {
		log.Debug("Serivce endpoint doesnt configure any argument schema")

	}
	r, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}
	err = m.isArgumentsValid(r, schema)
	if err != nil {
		return nil, err
	}
	req.Arguments = string(r)

	resp, err := svc.Call(ctx, req)
	if err != nil {
		log.Debug("Call return with error", "err", err.Error())
		return nil, err
	}
	if resp.Status == v1.Status_Error {
		log.Debug("Call return with error", "err", resp.Error)
		return []byte(resp.Payload), fmt.Errorf(resp.Error)
	}

	returnSchema, ok := schemas[fmt.Sprintf("endpoints/%s/%s", endpoint, "returns.json")]
	if !ok {
		log.Debug("Serivce endpoint doesnt configure any return schema")
		return []byte(resp.Payload), nil
	}
	err = m.isResponsePayloadValid(resp.Payload, returnSchema)
	if err != nil {
		return []byte(resp.Payload), err
	}

	return []byte(resp.Payload), nil
}

// Destory stop the modem and all the services
func (m *modem) Destroy() error {
	for name, service := range m.services {
		err := service.Kill()
		if err != nil {
			m.logger.Debug("Failed to kill service", "service", name)
		}
		m.logger.Debug("Service stopped", "service", name)
	}
	return nil
}

// AddService adds external service to the modem
func (m *modem) AddService(name string, runner service.Service) error {
	m.services[name] = runner
	return nil
}

func (m *modem) initService(name string, svc service.Service, log logger.Logger) error {
	if err := svc.Run(); err != nil {
		return err
	}
	return nil
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
