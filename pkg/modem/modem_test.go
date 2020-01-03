package modem

import (
	"errors"
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
	testingServicePort                    = "9090"
	testingErrorFailedToCreateLogFile     = "FailedToCreateLogFile"
	testingErrorFailedToDialToService     = "FailedToDialToService"
	testingErrorFailedToCallInitOnService = "FailedToCallInitOnService"
	testingErrorFailedRunService          = "FailedToRunService"
)

func buildMockService(runnerMockProbider func() *mocks.Runner) *service {
	svc := &service{
		id: testingServiceID,
		server: struct {
			port string
		}{
			port: testingServicePort,
		},
	}
	if runnerMockProbider != nil {
		svc.runner = runnerMockProbider()
	} else {
		r := &mocks.Runner{}
		r.On("Run", mock.Anything).Return(nil)
		svc.runner = r
	}
	return svc
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
			name: "Successfully initialized service, no error",
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
				serviceClientCreator: func(m *mockServiceClientCreator) *mockServiceClientCreator {
					client := mocks.ServiceClient{}
					client.On("Init", mock.Anything, mock.Anything).Return(&v1.InitResponse{}, nil)
					m.On("New", mock.Anything).Return(&client)
					return m
				},
				services: map[string]*service{
					"svc": buildMockService(nil),
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
					"svc": buildMockService(nil),
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
				services: map[string]*service{
					"svc": buildMockService(func() *mocks.Runner {
						m := &mocks.Runner{}
						m.On("Run", mock.Anything).Return(errors.New(testingErrorFailedRunService))
						return m
					}),
				},
			},
			wantErr: true,
			err:     "Serive: svc - Error: FailedToRunService\n",
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
				dialer: func(m *mockDialer) *mockDialer {
					m.On("Dial", "localhost:9090", mock.Anything).Return(nil, errors.New(testingErrorFailedToDialToService))
					return m
				},
				services: map[string]*service{
					"svc": buildMockService(nil),
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
				dialer: func(m *mockDialer) *mockDialer {
					m.On("Dial", "localhost:9090", mock.Anything).Return(nil, nil)
					return m
				},
				services: map[string]*service{
					"svc": buildMockService(nil),
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
