package core_test

import (
	"testing"

	"github.com/open-integration/core"
	apiv1 "github.com/open-integration/core/pkg/api/v1"
	"github.com/open-integration/core/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var pipelineTestMetadata = core.PipelineMetadata{
	Name: "pipeline",
}

func Test_engine_Run(t *testing.T) {
	tests := []struct {
		name       string
		options    *core.EngineOptions
		wantErr    bool
		middleware []func(e core.Engine)
	}{
		struct {
			name       string
			options    *core.EngineOptions
			wantErr    bool
			middleware []func(e core.Engine)
		}{
			name:    "Should run zero tasks with no errors",
			wantErr: false,
			options: &core.EngineOptions{
				Logger: extendLoggerMockWithBasicMocks(createFakeLogger()),
			},
		},
		struct {
			name       string
			options    *core.EngineOptions
			wantErr    bool
			middleware []func(e core.Engine)
		}{
			name:    "Should run one task once the engine started and exit succesfuly",
			wantErr: false,
			options: &core.EngineOptions{
				Logger: extendLoggerMockWithBasicMocks(createFakeLogger()),
				Pipeline: core.Pipeline{
					Metadata: pipelineTestMetadata,
					Spec: core.PipelineSpec{
						Services: []core.Service{
							core.Service{
								Name:    "some-service",
								As:      "service-name",
								Version: "0.0.1",
							},
						},
						Reactions: []core.EventReaction{
							core.EventReaction{
								Condition: core.ConditionEngineStarted,
								Reaction: func(ev core.Event, state core.State) []core.Task {
									return []core.Task{
										core.Task{
											Metadata: core.TaskMetadata{
												Name: "task-name",
											},
											Spec: core.TaskSpec{
												Endpoint: "endpoint",
												Service:  "service-name",
											},
										},
									}
								},
							},
						},
					},
				},
			},
			middleware: []func(e core.Engine){
				func(e core.Engine) {
					runner := extendRunnerMockWithBasicMocks(createMockedRunner())
					e.Modem().AddService("service-id", "service-name", runner)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := core.NewEngine(tt.options)
			for _, m := range tt.middleware {
				m(e)
			}
			if tt.wantErr {
				assert.Error(t, e.Run())
			} else {
				assert.NoError(t, e.Run())
			}
		})
	}
}

func createMockedRunner() *mocks.Runner {
	return &mocks.Runner{}
}

func extendRunnerMockWithBasicMocks(m *mocks.Runner) *mocks.Runner {
	m.On("Run").Return(nil)
	m.On("Kill").Return(nil)
	m.On("Call", mock.Anything, mock.Anything).Return(&apiv1.CallResponse{
		Status: apiv1.Status_OK,
	}, nil)
	return m
}

func createFakeLogger() *mocks.Logger {
	l := &mocks.Logger{}
	return l
}

func extendLoggerMockWithBasicMocks(m *mocks.Logger) *mocks.Logger {
	m.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("New", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
	return m
}
