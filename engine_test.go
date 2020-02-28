package core_test

import (
	"testing"

	"github.com/open-integration/core"
	apiv1 "github.com/open-integration/core/pkg/api/v1"
	"github.com/open-integration/core/pkg/mocks"
	"github.com/open-integration/core/pkg/state"
	"github.com/open-integration/core/pkg/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var pipelineTestMetadata = core.PipelineMetadata{
	Name: "pipeline",
}

type (
	testEngineRun struct {
		name       string
		options    *core.EngineOptions
		wantErr    bool
		middleware []func(e core.Engine)
	}
)

func Test_engine_Run(t *testing.T) {
	tests := []testEngineRun{
		testEngineRun{
			name:    "Should run zero tasks with no errors",
			wantErr: false,
			options: &core.EngineOptions{
				Logger: extendLoggerMockWithBasicMocks(createFakeLogger()),
			},
		},
		testEngineRun{
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
								Reaction: func(ev state.Event, state state.State) []task.Task {
									return []task.Task{
										task.Task{
											Metadata: task.Metadata{
												Name: "task-name",
											},
											Spec: task.Spec{
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
		testEngineRun{
			name:    "Should create multiple tasks as a result of previous task",
			wantErr: false,
			options: &core.EngineOptions{
				Logger: extendLoggerMockWithBasicMocks(createFakeLogger()),
				Pipeline: core.Pipeline{
					Metadata: pipelineTestMetadata,
					Spec: core.PipelineSpec{
						Services: []core.Service{
							core.Service{
								Name:    "service(task1)",
								As:      "service1",
								Version: "0.0.1",
							},
							core.Service{
								Name:    "service(task2...5)",
								As:      "service2",
								Version: "0.0.1",
							},
						},
						Reactions: []core.EventReaction{
							core.EventReaction{
								Condition: core.ConditionEngineStarted,
								Reaction: func(ev state.Event, state state.State) []task.Task {
									return []task.Task{
										task.Task{
											Metadata: task.Metadata{
												Name: "task-1",
											},
											Spec: task.Spec{
												Service:  "service1",
												Endpoint: "endpoint1",
											},
										},
									}
								},
							},
							core.EventReaction{
								Condition: core.ConditionTaskFinishedWithStatus("task-1", state.TaskStatusSuccess),
								Reaction: func(ev state.Event, state state.State) []task.Task {
									return []task.Task{
										task.Task{
											Metadata: task.Metadata{
												Name: "task-2",
											},
											Spec: task.Spec{
												Service:  "service2",
												Endpoint: "endpoint1",
											},
										},
										task.Task{
											Metadata: task.Metadata{
												Name: "task-3",
											},
											Spec: task.Spec{
												Service:  "service2",
												Endpoint: "endpoint1",
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
					e.Modem().AddService("service-id-1", "service1", runner)
				},
				func(e core.Engine) {
					runner := extendRunnerMockWithBasicMocks(createMockedRunner())
					e.Modem().AddService("service-id-2", "service2", runner)
				},
			},
		},
		testEngineRun{
			name:    "Should run succesfully multiple tasks",
			wantErr: false,
			middleware: []func(e core.Engine){
				func(e core.Engine) {
					runner := extendRunnerMockWithBasicMocks(createMockedRunner())
					e.Modem().AddService("service-id-1", "service1", runner)
				},
				func(e core.Engine) {
					runner := extendRunnerMockWithBasicMocks(createMockedRunner())
					e.Modem().AddService("service-id-2", "service2", runner)
				},
			},
			options: &core.EngineOptions{
				Logger: extendLoggerMockWithBasicMocks(createFakeLogger()),
				Pipeline: core.Pipeline{
					Metadata: pipelineTestMetadata,
					Spec: core.PipelineSpec{
						Services: []core.Service{
							core.Service{
								Name:    "service(task1)",
								As:      "service1",
								Version: "0.0.1",
							},
							core.Service{
								Name:    "service(task2...5)",
								As:      "service2",
								Version: "0.0.1",
							},
						},
						Reactions: []core.EventReaction{
							core.EventReaction{
								Condition: core.ConditionEngineStarted,
								Reaction: func(ev state.Event, state state.State) []task.Task {
									return []task.Task{
										task.Task{
											Metadata: task.Metadata{
												Name: "task-1",
											},
											Spec: task.Spec{
												Service:  "service1",
												Endpoint: "endpoint1",
											},
										},
									}
								},
							},
							core.EventReaction{
								Condition: core.ConditionTaskFinishedWithStatus("task-1", state.TaskStatusSuccess),
								Reaction: func(ev state.Event, state state.State) []task.Task {
									return []task.Task{
										task.Task{
											Metadata: task.Metadata{
												Name: "task-2",
											},
											Spec: task.Spec{
												Service:  "service2",
												Endpoint: "endpoint1",
											},
										},
										task.Task{
											Metadata: task.Metadata{
												Name: "task-3",
											},
											Spec: task.Spec{
												Service:  "service2",
												Endpoint: "endpoint1",
											},
										},
									}
								},
							},
							core.EventReaction{
								Condition: core.ConditionTaskFinishedWithStatus("task-2", state.TaskStatusSuccess),
								Reaction: func(ev state.Event, state state.State) []task.Task {
									return []task.Task{
										task.Task{
											Metadata: task.Metadata{
												Name: "task-4",
											},
											Spec: task.Spec{
												Service:  "service2",
												Endpoint: "endpoint1",
											},
										},
										task.Task{
											Metadata: task.Metadata{
												Name: "task-5",
											},
											Spec: task.Spec{
												Service:  "service2",
												Endpoint: "endpoint1",
											},
										},
									}
								},
							},
						},
					},
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
