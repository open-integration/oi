package templates

// PipelineTemplate template
var PipelineTemplate = `
package main

import (
	"fmt"
	"context"

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
							oi.NewFunctionTask("{{ .Name  }}", Run{{ .Name }}),
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


func Run{{ .Name }}(ctx context.Context, options task.RunOptions) ([]byte, error) {
	fmt.Println("*****************")
	fmt.Println("Hello World")
	fmt.Println("*****************")
	return []byte{}, nil
}

{{ end }}
{{ end }}`
