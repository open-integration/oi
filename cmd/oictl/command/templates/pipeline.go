package templates

// PipelineTemplate template
var PipelineTemplate = `
package main

import (
	"fmt"
	"context"
	"strings"
	"os/exec"

	"github.com/open-integration/oi"
	"github.com/open-integration/oi/core/engine"
	"github.com/open-integration/oi/core/event"
	"github.com/open-integration/oi/core/state"
	"github.com/open-integration/oi/core/task"
)

func main() {
	pipe := engine.Pipeline{
		Metadata: engine.PipelineMetadata{
			Name: "{{ .Metadata.Name }}",
		},
		Spec: engine.PipelineSpec{
			Services: []engine.Service{
			{{ range .Services }}
				engine.Service{
					As:      "{{ .As }}",
					Name:    "{{ .Name }}",
					Version: "{{ .Version }}",
				},
			{{ end }}
			},
			Reactions: []engine.EventReaction{
				{{ range $i, $eventReaction := .EventReactions -}}
				{
					Condition: {{ .Condition }},
					Reaction: func(ev event.Event, state state.State) []task.Task {
						return []task.Task{
							{{ range $eventReaction.Reaction -}}
							oi.NewFunctionTask("{{ .Name  }}", Run{{ .Name }}("{{ .Command }}")),
							{{ end }}
						}
					},
				},
				{{- end }}
			},
		},
	}
	e := oi.NewEngine(&oi.EngineOptions{
		Pipeline: pipe,
	})
	oi.HandleEngineError(e.Run())
}

{{ range $i, $eventReaction := .EventReactions }}
{{ range $eventReaction.Reaction -}}

func Run{{ .Name }}(command string) func(ctx context.Context, opt task.RunOptions) ([]byte, error) {
	return func(ctx context.Context, opt task.RunOptions) ([]byte, error) {
		c := strings.Split(command, " ")
		res := exec.CommandContext(ctx, c[0], strings.Join(c[1:], ""))
		data, err := res.Output()
		if err != nil {
			return nil, err
		}
		fmt.Println(string(data))          // Print to the stdout
		// fmt.Fprintln(opt.FD, string(data)) // Print to task logger, later the file can be found in $PWD/logs/tasks/{{ .Name }}.log
		return []byte{}, nil
	}
}

{{ end }}
{{ end }}`
