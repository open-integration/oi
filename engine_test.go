package core

import (
	"testing"

	apiv1 "github.com/open-integration/core/pkg/api/v1"
	"github.com/open-integration/core/pkg/event"
	"github.com/open-integration/core/pkg/mocks"
	"github.com/open-integration/core/pkg/service"
	"github.com/open-integration/core/pkg/state"
	"github.com/open-integration/core/pkg/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var pipelineTestMetadata = PipelineMetadata{
	Name: "pipeline",
}

type (
	testEngineRun struct {
		name       string
		options    *EngineOptions
		wantErr    bool
		middleware []func(e Engine)
	}

	fakeService struct {
		name   string
		runner service.Service
	}
)

func Test_engine_Run(t *testing.T) {
	tests := []testEngineRun{
		{
			name:    "Should run zero tasks with no errors",
			wantErr: false,
			options: &EngineOptions{
				Logger: extendLoggerMockWithBasicMocks(createFakeLogger()),
			},
		},
		{
			name:    "Should run one task once the engine started and exit succesfuly",
			wantErr: false,
			options: &EngineOptions{
				Logger:            extendLoggerMockWithBasicMocks(createFakeLogger()),
				serviceDownloader: createFakeDownloader(),
				modem: createFakeModem([]fakeService{
					{
						name:   "some-service",
						runner: createFakeServiceRunner(),
					},
				}),
				Pipeline: Pipeline{
					Metadata: pipelineTestMetadata,
					Spec: PipelineSpec{
						Services: []Service{
							{
								Name:    "some-service",
								As:      "service-name",
								Version: "0.0.1",
							},
						},
						Reactions: []EventReaction{
							{
								Condition: ConditionEngineStarted(),
								Reaction: func(ev event.Event, state state.State) []task.Task {
									return []task.Task{
										{
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
			middleware: []func(e Engine){
				func(e Engine) {
					runner := extendRunnerMockWithBasicMocks(createFakeServiceRunner())
					e.Modem().AddService("service-name", runner)
				},
			},
		},
		{
			name:    "Should create multiple tasks as a result of previous task",
			wantErr: false,

			options: &EngineOptions{
				Logger:            extendLoggerMockWithBasicMocks(createFakeLogger()),
				serviceDownloader: createFakeDownloader(),
				modem: createFakeModem([]fakeService{
					{
						name:   "some-service",
						runner: createFakeServiceRunner(),
					},
				}),
				Pipeline: Pipeline{
					Metadata: pipelineTestMetadata,
					Spec: PipelineSpec{
						Services: []Service{
							{
								Name:    "service(task1)",
								As:      "service1",
								Version: "0.0.1",
							},
							{
								Name:    "service(task2...5)",
								As:      "service2",
								Version: "0.0.1",
							},
						},
						Reactions: []EventReaction{
							{
								Condition: ConditionEngineStarted(),
								Reaction: func(ev event.Event, state state.State) []task.Task {
									return []task.Task{
										{
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
							{
								Condition: ConditionTaskFinishedWithStatus("task-1", state.TaskStatusSuccess),
								Reaction: func(ev event.Event, state state.State) []task.Task {
									return []task.Task{
										{
											Metadata: task.Metadata{
												Name: "task-2",
											},
											Spec: task.Spec{
												Service:  "service2",
												Endpoint: "endpoint1",
											},
										},
										{
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
			middleware: []func(e Engine){
				func(e Engine) {
					runner := extendRunnerMockWithBasicMocks(createFakeServiceRunner())
					e.Modem().AddService("service1", runner)
				},
				func(e Engine) {
					runner := extendRunnerMockWithBasicMocks(createFakeServiceRunner())
					e.Modem().AddService("service2", runner)
				},
			},
		},
		{
			name:    "Should run succesfully multiple tasks",
			wantErr: false,
			middleware: []func(e Engine){
				func(e Engine) {
					runner := extendRunnerMockWithBasicMocks(createFakeServiceRunner())
					e.Modem().AddService("service1", runner)
				},
				func(e Engine) {
					runner := extendRunnerMockWithBasicMocks(createFakeServiceRunner())
					e.Modem().AddService("service2", runner)
				},
			},
			options: &EngineOptions{
				Logger:            extendLoggerMockWithBasicMocks(createFakeLogger()),
				serviceDownloader: createFakeDownloader(),
				modem: createFakeModem([]fakeService{
					{
						name:   "some-service",
						runner: createFakeServiceRunner(),
					},
				}),
				Pipeline: Pipeline{
					Metadata: pipelineTestMetadata,
					Spec: PipelineSpec{
						Services: []Service{
							{
								Name:    "service(task1)",
								As:      "service1",
								Version: "0.0.1",
							},
							{
								Name:    "service(task2...5)",
								As:      "service2",
								Version: "0.0.1",
							},
						},
						Reactions: []EventReaction{
							{
								Condition: ConditionEngineStarted(),
								Reaction: func(ev event.Event, state state.State) []task.Task {
									return []task.Task{
										{
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
							{
								Condition: ConditionTaskFinishedWithStatus("task-1", state.TaskStatusSuccess),
								Reaction: func(ev event.Event, state state.State) []task.Task {
									return []task.Task{
										{
											Metadata: task.Metadata{
												Name: "task-2",
											},
											Spec: task.Spec{
												Service:  "service2",
												Endpoint: "endpoint1",
											},
										},
										{
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
							{
								Condition: ConditionTaskFinishedWithStatus("task-2", state.TaskStatusSuccess),
								Reaction: func(ev event.Event, state state.State) []task.Task {
									return []task.Task{
										{
											Metadata: task.Metadata{
												Name: "task-4",
											},
											Spec: task.Spec{
												Service:  "service2",
												Endpoint: "endpoint1",
											},
										},
										{
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
			e := NewEngine(tt.options)
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

func extendRunnerMockWithBasicMocks(m *mocks.Service) *mocks.Service {
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

func createFakeDownloader() *mocks.Downloader {
	d := &mocks.Downloader{}
	d.On("Download", mock.Anything, mock.Anything).Return("", nil)
	return d
}

func createFakeServiceRunner() *mocks.Service {
	r := &mocks.Service{}
	r.On("Schemas").Return(map[string]string{})
	return r
}

func createFakeModem(services []fakeService) *mocks.Modem {
	m := &mocks.Modem{}
	m.On("AddService", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]byte(""), nil)
	m.On("Init").Return(nil)
	m.On("Destroy").Return(nil)
	for _, s := range services {
		m.AddService(s.name, s.runner)
	}
	return m
}
