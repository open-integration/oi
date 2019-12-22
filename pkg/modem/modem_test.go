package modem

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	v1 "github.com/open-integration/core/pkg/api/v1"
	"github.com/open-integration/core/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testingServiceDirectory               = "testing-service-directory"
	testingServiceID                      = "service-id"
	testingServiceBinaryPath              = "/path/to/bin"
	testingServicePort                    = "9090"
	testingServicePid                     = 1
	testingErrorFailedToCreateLogFile     = "FailedToCreateLogFile"
	testingErrorFailedToStartService      = "FailedToStartService"
	testingErrorFailedToDialToService     = "FailedToDialToService"
	testingErrorFailedToCallInitOnService = "FailedToCallInitOnService"
)

func buildMockService() *service {
	return &service{
		id: testingServiceID,
		server: struct {
			binPath string
			port    string
			pid     int
		}{
			binPath: testingServiceBinaryPath,
			port:    testingServicePort,
			pid:     testingServicePid,
		},
	}
}

func buildBasicLoggerMock(m *mocks.Logger) *mocks.Logger {
	m.On("Debug", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("New", mock.Anything, mock.Anything).Return(m)
	return m
}

func Test_modem_Init(t *testing.T) {
	type fields struct {
		services             map[string]*service
		logger               func(*mocks.Logger) *mocks.Logger
		serviceLogDirectory  string
		wg                   *sync.WaitGroup
		logFileCreator       func(m *mockFileCreator) *mockFileCreator
		dialer               func(m *mockDialer) *mockDialer
		serviceStarter       func(m *mockServiceStarter) *mockServiceStarter
		serviceClientCreator func(m *mockServiceClientCreator) *mockServiceClientCreator
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		err     string
	}{
		{
			name: "No service, not possilbe error",
			fields: fields{
				wg:     &sync.WaitGroup{},
				logger: buildBasicLoggerMock,
			},
			wantErr: false,
		},
		{
			name: "Succesfully initialized service, no error",
			fields: fields{
				wg:                  &sync.WaitGroup{},
				serviceLogDirectory: testingServiceDirectory,
				logger:              buildBasicLoggerMock,
				dialer: func(m *mockDialer) *mockDialer {
					m.On("Dial", "localhost:9090", mock.Anything).Return(nil, nil)
					return m
				},
				logFileCreator: func(m *mockFileCreator) *mockFileCreator {
					m.On("Create", testingServiceDirectory, "svc-service-id.log").Return(nil, nil)
					return m
				},
				serviceStarter: func(m *mockServiceStarter) *mockServiceStarter {
					m.On("Exec", testingServiceBinaryPath, []string{""}, []string{fmt.Sprintf("PORT=%s", testingServicePort)}, true, "", "", nil).Return(0, nil)
					return m
				},
				serviceClientCreator: func(m *mockServiceClientCreator) *mockServiceClientCreator {
					client := mocks.ServiceClient{}
					client.On("Init", mock.Anything, mock.Anything).Return(&v1.InitResponse{}, nil)
					m.On("New", mock.Anything).Return(&client)
					return m
				},
				services: map[string]*service{
					"svc": buildMockService(),
				},
			},
			wantErr: false,
		},
		{
			name: "Failed to create log file, exit with error",
			fields: fields{
				wg:                  &sync.WaitGroup{},
				serviceLogDirectory: testingServiceDirectory,
				logger:              buildBasicLoggerMock,
				logFileCreator: func(m *mockFileCreator) *mockFileCreator {
					m.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New(testingErrorFailedToCreateLogFile))
					return m
				},
				services: map[string]*service{
					"svc": buildMockService(),
				},
			},
			wantErr: true,
			err:     "Serive: svc - Error: FailedToCreateLogFile\n",
		},
		{
			name: "Failed to start service, exit with error",
			fields: fields{
				wg:                  &sync.WaitGroup{},
				serviceLogDirectory: testingServiceDirectory,
				logger:              buildBasicLoggerMock,
				logFileCreator: func(m *mockFileCreator) *mockFileCreator {
					m.On("Create", testingServiceDirectory, "svc-service-id.log").Return(nil, nil)
					return m
				},
				serviceStarter: func(m *mockServiceStarter) *mockServiceStarter {
					m.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New(testingErrorFailedToStartService))
					return m
				},
				services: map[string]*service{
					"svc": buildMockService(),
				},
			},
			wantErr: true,
			err:     "Serive: svc - Error: FailedToStartService\n",
		},
		{
			name: "Failed to dial service, exit with error",
			fields: fields{
				wg:                  &sync.WaitGroup{},
				serviceLogDirectory: testingServiceDirectory,
				logger:              buildBasicLoggerMock,
				logFileCreator: func(m *mockFileCreator) *mockFileCreator {
					m.On("Create", testingServiceDirectory, "svc-service-id.log").Return(nil, nil)
					return m
				},
				serviceStarter: func(m *mockServiceStarter) *mockServiceStarter {
					m.On("Exec", testingServiceBinaryPath, []string{""}, []string{fmt.Sprintf("PORT=%s", testingServicePort)}, true, "", "", nil).Return(0, nil)
					return m
				},
				dialer: func(m *mockDialer) *mockDialer {
					m.On("Dial", "localhost:9090", mock.Anything).Return(nil, errors.New(testingErrorFailedToDialToService))
					return m
				},
				services: map[string]*service{
					"svc": buildMockService(),
				},
			},
			wantErr: true,
			err:     "Serive: svc - Error: FailedToDialToService\n",
		},
		{
			name: "Failed to call init, exit with error",
			fields: fields{
				wg:                  &sync.WaitGroup{},
				serviceLogDirectory: testingServiceDirectory,
				logger:              buildBasicLoggerMock,
				logFileCreator: func(m *mockFileCreator) *mockFileCreator {
					m.On("Create", testingServiceDirectory, "svc-service-id.log").Return(nil, nil)
					return m
				},
				serviceStarter: func(m *mockServiceStarter) *mockServiceStarter {
					m.On("Exec", testingServiceBinaryPath, []string{""}, []string{fmt.Sprintf("PORT=%s", testingServicePort)}, true, "", "", nil).Return(0, nil)
					return m
				},
				dialer: func(m *mockDialer) *mockDialer {
					m.On("Dial", "localhost:9090", mock.Anything).Return(nil, nil)
					return m
				},
				services: map[string]*service{
					"svc": buildMockService(),
				},
				serviceClientCreator: func(m *mockServiceClientCreator) *mockServiceClientCreator {
					client := mocks.ServiceClient{}
					client.On("Init", mock.Anything, mock.Anything).Return(nil, errors.New(testingErrorFailedToCallInitOnService))
					m.On("New", mock.Anything).Return(&client)
					return m
				},
			},
			wantErr: true,
			err:     "Serive: svc - Error: FailedToCallInitOnService\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := &modem{
				services:            tt.fields.services,
				serviceLogDirectory: tt.fields.serviceLogDirectory,
				wg:                  tt.fields.wg,
			}
			if tt.fields.logger != nil {
				m.logger = tt.fields.logger(&mocks.Logger{})
			}
			if tt.fields.logFileCreator != nil {
				m.logFileCreator = tt.fields.logFileCreator(&mockFileCreator{})
			}
			if tt.fields.dialer != nil {
				m.dialer = tt.fields.dialer(&mockDialer{})
			}
			if tt.fields.serviceStarter != nil {
				m.serviceStarter = tt.fields.serviceStarter(&mockServiceStarter{})
			}
			if tt.fields.serviceClientCreator != nil {
				m.serviceClientCreator = tt.fields.serviceClientCreator(&mockServiceClientCreator{})
			}
			err := m.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("modem.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				assert.Equal(t, err.Error(), tt.err)
			}
		})
	}
}
