package core_test

import (
	"testing"

	"github.com/open-integration/core"
	"github.com/open-integration/core/pkg/state"
	"github.com/stretchr/testify/assert"
)

func TestConditionTaskFinishedWithStatus(t *testing.T) {
	type args struct {
		task   string
		status string
		ev     state.Event
		state  func() state.State
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		struct {
			name string
			args args
			want bool
		}{
			name: "Should return true when task finished with the same status",
			want: true,
			args: args{
				status: state.EventTaskFinished,
				task:   "task-1",
				ev: state.Event{
					Metadata: state.EventMetadata{
						Name: state.EventTaskFinished,
					},
				},
				state: func() state.State {
					return nil
					// Tasks: map[core.ID]state.TaskState{
					// 	"task-id-1": state.TaskState{
					// 		Status: state.EventTaskFinished,
					// 		State:  "finished",
					// 		Task:   "task-1",
					// 	},
					// },
				},
			},
		},
		{
			name: "Should return false when the passed event not matched to the reuired event",
			want: false,
			args: args{
				status: state.EventTaskFinished,
				task:   "task-1",
				ev: state.Event{
					Metadata: state.EventMetadata{
						Name: state.EventTaskStarted,
					},
				},
			},
		},
		{
			name: "Should return false in case task was not finished or found in state task list",
			want: false,
			args: args{
				status: state.EventTaskStarted,
				ev: state.Event{
					Metadata: state.EventMetadata{
						Name: state.EventTaskFinished,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fn := core.ConditionTaskFinishedWithStatus(tt.args.task, tt.args.status)
			// assert.Equal(t, tt.want, fn(tt.args.ev, tt.args.state))
		})
	}
}

func TestConditionEngineStarted(t *testing.T) {
	type args struct {
		ev    state.Event
		state state.State
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		struct {
			name string
			args args
			want bool
		}{
			name: "Should return true in case the event was engine.started",
			want: true,
			args: args{
				ev: state.Event{
					Metadata: state.EventMetadata{
						Name: state.EventEngineStarted,
					},
				},
			},
		},
		struct {
			name string
			args args
			want bool
		}{
			name: "Should return false in case the event was not engine.started",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, core.ConditionEngineStarted(tt.args.ev, tt.args.state))
		})
	}
}
