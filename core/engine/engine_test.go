package engine

import (
	"testing"

	"github.com/open-integration/core/core/condition"
	"github.com/open-integration/core/core/event"
	"github.com/open-integration/core/core/modem"
	"github.com/open-integration/core/core/service/runner"
	"github.com/open-integration/core/core/state"
	"github.com/open-integration/core/core/task"
	"github.com/open-integration/core/core/task/tasks"
	apiv1 "github.com/open-integration/core/pkg/api/v1"
	"github.com/open-integration/core/pkg/downloader"
	"github.com/open-integration/core/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var pipelineTestMetadata = PipelineMetadata{
	Name: "pipeline",
}

type (
	testEngineRun struct {
		name       string
		options    *Options
		wantErr    bool
		middleware []func(e Engine)
	}

	fakeService struct {
		name   string
		runner runner.Service
	}
)

func Test_engine_Run(t *testing.T) {
	tests := []testEngineRun{
		{
			name:    "Should run zero tasks with no errors",
			wantErr: false,
			options: &Options{
				Logger: extendLoggerMockWithBasicMocks(createFakeLogger()),
			},
		},
		{
			name:    "Should run one task once the engine started and exit successfully",
			wantErr: false,
			options: &Options{
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
								Condition: condition.EngineStarted(),
								Reaction: func(ev event.Event, state state.State) []task.Task {
									return []task.Task{
										tasks.NewSerivceTask("task-name", "service-name", "endpoint"),
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

			options: &Options{
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
								Condition: condition.EngineStarted(),
								Reaction: func(ev event.Event, state state.State) []task.Task {
									return []task.Task{
										tasks.NewSerivceTask("task-1", "service1", "endpoint1"),
									}
								},
							},
							{
								Condition: condition.TaskFinishedWithStatus("task-1", state.TaskStatusSuccess),
								Reaction: func(ev event.Event, state state.State) []task.Task {
									return []task.Task{
										tasks.NewSerivceTask("task-2", "service2", "endpoint1"),
										tasks.NewSerivceTask("task-3", "service2", "endpoint1"),
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
			name:    "Should run successfully multiple tasks",
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
			options: &Options{
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
								Condition: condition.EngineStarted(),
								Reaction: func(ev event.Event, state state.State) []task.Task {
									return []task.Task{
										tasks.NewSerivceTask("task-1", "service1", "endpoint1"),
									}
								},
							},
							{
								Condition: condition.TaskFinishedWithStatus("task-1", state.TaskStatusSuccess),
								Reaction: func(ev event.Event, state state.State) []task.Task {
									return []task.Task{
										tasks.NewSerivceTask("task-2", "service2", "endpoint1"),
										tasks.NewSerivceTask("task-3", "service2", "endpoint1"),
									}
								},
							},
							{
								Condition: condition.TaskFinishedWithStatus("task-2", state.TaskStatusSuccess),
								Reaction: func(ev event.Event, state state.State) []task.Task {
									return []task.Task{
										tasks.NewSerivceTask("task-4", "service2", "endpoint1"),
										tasks.NewSerivceTask("task-5", "service2", "endpoint1"),
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
			e, err := New(tt.options)
			assert.NoError(t, err)
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

func extendRunnerMockWithBasicMocks(m *runner.MockService) *runner.MockService {
	m.On("Run").Return(nil)
	m.On("Kill").Return(nil)
	m.
		On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&apiv1.CallResponse{
			Status: apiv1.Status_OK,
		}, nil)
	return m
}

func createFakeLogger() *logger.MockLogger {
	l := &logger.MockLogger{}
	return l
}

func extendLoggerMockWithBasicMocks(m *logger.MockLogger) logger.Logger {
	m.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("New", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
	return m
}

func createFakeDownloader() downloader.Downloader {
	d := &downloader.MockDownloader{}
	d.On("Download", mock.Anything, mock.Anything).Return("", nil)
	return d
}

func createFakeServiceRunner() *runner.MockService {
	r := &runner.MockService{}
	r.On("Schemas").Return(map[string]string{})
	return r
}

func createFakeModem(services []fakeService) modem.Modem {
	m := &modem.MockModem{}
	m.On("AddService", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("Call", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]byte(""), nil)
	m.On("Init").Return(nil)
	m.On("Destroy").Return(nil)
	for _, s := range services {
		m.AddService(s.name, s.runner)
	}
	return m
}
