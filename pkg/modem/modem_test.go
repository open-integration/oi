package modem

import (
	"errors"
	"sync"
	"testing"

	"github.com/open-integration/core/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testingServiceDirectory      = "testing-service-directory"
	testingServiceID             = "service-id"
	testingErrorFailedRunService = "FailedToRunService"
)

func buildMockService(runnerMockProvider func() *mocks.Runner) *service {
	svc := &service{
		id: testingServiceID,
	}
	if runnerMockProvider != nil {
		svc.runner = runnerMockProvider()
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
		services            map[string]*service
		logger              func(*mocks.Logger) *mocks.Logger
		serviceLogDirectory string
		wg                  *sync.WaitGroup
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
				services: map[string]*service{
					"svc": buildMockService(nil),
				},
			},
			wantErr: false,
		},
		{
			name: "Failed to start service, exit with error",
			fields: fields{
				wg:                  &sync.WaitGroup{},
				serviceLogDirectory: testingServiceDirectory,
				logger:              buildBasicLoggerMock,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := &modem{
				services: tt.fields.services,
				wg:       tt.fields.wg,
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
