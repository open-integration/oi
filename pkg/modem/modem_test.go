package modem

import (
	"errors"
	"testing"

	v1 "github.com/open-integration/core/pkg/api/v1"
	"github.com/open-integration/core/pkg/logger"
	"github.com/open-integration/core/pkg/mocks"
	"github.com/open-integration/core/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testingServiceDirectory      = "testing-service-directory"
	testingServiceID             = "service-id"
	testingErrorFailedRunService = "FailedToRunService"
)

func buildMockService(runnerMockProvider func() *mocks.Service) service.Service {
	if runnerMockProvider != nil {
		return runnerMockProvider()
	}
	r := &mocks.Service{}
	r.On("Run", mock.Anything).Return(nil)
	r.On("Schemas").Return(map[string]string{})
	return r
}

func buildBasicLoggerMock(m *mocks.Logger) *mocks.Logger {
	m.On("Debug", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("New", mock.Anything, mock.Anything).Return(m)
	m.On("New", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
	return m
}

func Test_modem_Init(t *testing.T) {
	type fields struct {
		services            map[string]service.Service
		logger              func(*mocks.Logger) *mocks.Logger
		serviceLogDirectory string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		err     string
	}{
		{
			name: "No service, no possilbe error",
			fields: fields{
				logger: buildBasicLoggerMock,
			},
			wantErr: false,
		},
		{
			name: "Successfully initialized service, no error",
			fields: fields{
				serviceLogDirectory: testingServiceDirectory,
				logger:              buildBasicLoggerMock,
				services: map[string]service.Service{
					"svc": buildMockService(nil),
				},
			},
			wantErr: false,
		},
		{
			name: "Failed to start service, exit with error",
			fields: fields{
				serviceLogDirectory: testingServiceDirectory,
				logger:              buildBasicLoggerMock,
				services: map[string]service.Service{
					"svc": buildMockService(func() *mocks.Service {
						m := &mocks.Service{}
						m.On("Schemas").Return(map[string]string{})
						m.On("Run", mock.Anything).Return(errors.New(testingErrorFailedRunService))
						return m
					}),
				},
			},
			wantErr: true,
			err:     "Failed to initiate service: FailedToRunService",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := &modem{
				services: tt.fields.services,
			}
			if tt.fields.logger != nil {
				m.logger = tt.fields.logger(&mocks.Logger{})
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

func Test_modem_Call(t *testing.T) {
	type fields struct {
		logger   logger.Logger
		services map[string]service.Service
	}
	type args struct {
		service   string
		endpoint  string
		arguments map[string]interface{}
		fd        string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Should return an error in case the service was not found",
			args: args{
				service: "not-found",
			},
			wantErr: true,
			fields: fields{
				logger: buildBasicLoggerMock(&mocks.Logger{}),
			},
		}, {
			name: "Should call the service with empty arguments in case it doesnt have argument schema",
			args: args{
				endpoint:  "endpoint",
				service:   "service",
				arguments: map[string]interface{}{},
			},
			fields: fields{
				services: map[string]service.Service{
					"service": buildMockService(func() *mocks.Service {
						m := &mocks.Service{}
						m.On("Schemas").Return(map[string]string{})
						m.On("Call", mock.Anything, &v1.CallRequest{
							Endpoint:  "endpoint",
							Arguments: "{}",
						}).Return(&v1.CallResponse{}, nil)
						return m
					}),
				},
				logger: buildBasicLoggerMock(&mocks.Logger{}),
			},
			want: []byte{},
		}, {
			name: "Should call the service with the arguments in case they mached to the schema",
			args: args{
				endpoint: "endpoint",
				service:  "service",
				arguments: map[string]interface{}{
					"key": "value",
				},
			},
			fields: fields{
				services: map[string]service.Service{
					"service": func() service.Service {
						s := buildMockService(func() *mocks.Service {
							m := &mocks.Service{}
							m.On("Schemas").Return(map[string]string{
								"endpoint/arguments.json": "{\"properties\": { \"key\": { \"type\": \"string\" } } }",
							})
							m.On("Call", mock.Anything, &v1.CallRequest{
								Endpoint:  "endpoint",
								Arguments: "{\"key\":\"value\"}",
							}).Return(&v1.CallResponse{}, nil)
							return m
						})
						return s
					}(),
				},
				logger: buildBasicLoggerMock(&mocks.Logger{}),
			},
			want: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &modem{
				services: tt.fields.services,
				logger:   tt.fields.logger,
			}
			got, err := m.Call(tt.args.service, tt.args.endpoint, tt.args.arguments, tt.args.fd)
			if (err != nil) != tt.wantErr {
				t.Errorf("modem.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
