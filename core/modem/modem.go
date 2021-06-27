package modem

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/open-integration/oi/core/service/runner"
	v1 "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
)

type (
	// Modem holds all the connection the external services
	Modem interface {
		ServiceCaller
		Init() error
		Destroy() error
		AddService(name string, runner runner.Service) error
	}

	modem struct {
		services map[string]runner.Service
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
		services: map[string]runner.Service{},
	}
	return m
}

// Init starts the modem and all the services
func (m *modem) Init() error {
	m.logger.Debug("Modem initiation started")
	for name, s := range m.services {
		if err := m.initService(name, s, m.logger.New("service", name)); err != nil {
			return fmt.Errorf("failed to initiate service: %w", err)
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
		return nil, fmt.Errorf("service %s not found", service)
	}
	r, err := json.Marshal(arguments)
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

	return []byte(resp.Payload), nil
}

// Destroy stop the modem and all the services
func (m *modem) Destroy() error {
	for name, service := range m.services {
		err := service.Kill()
		if err != nil {
			m.logger.Debug("failed to kill service", "service", name)
		}
		m.logger.Debug("Service stopped", "service", name)
	}
	return nil
}

// AddService adds external service to the modem
func (m *modem) AddService(name string, runner runner.Service) error {
	m.services[name] = runner
	return nil
}

func (m *modem) initService(name string, svc runner.Service, log logger.Logger) error {
	if err := svc.Run(); err != nil {
		return err
	}
	return nil
}
