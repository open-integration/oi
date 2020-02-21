# Open Integration - a pipeline execution engine

[![Go Report Card](https://goreportcard.com/badge/github.com/open-integration/core)](https://goreportcard.com/report/github.com/open-integration/core)
[![codecov](https://codecov.io/gh/open-integration/core/branch/master/graph/badge.svg)](https://codecov.io/gh/open-integration/core)
[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/olegs-codefresh/open-integration%2Fcore?type=cf-1)]( https%3A%2F%2Fg.codefresh.io%2Fpublic%2Faccounts%2Folegs-codefresh%2Fpipelines%2F5df37658c4bb05f822229465)

## stop using YAMLs

## This is not production ready stuff
Until the project has not reached version > 1.x.x it may have breaking changes in the API, please use the latest version before opening issue.

## Concepts
* A compiled, binray pipeline
* State - the engine holds the state of all the tasks
* Service - a standalone binary exposing API over http2 (gRPC) that the engine can trigger, called endpoint.
* Task - a request from a service to run some logic.

* Endpoint of a service defined by 2 files of JSON schema, `arguments.json` and `returns.json`, the engine will enforce the arguments given by a task and the output created to match the schema.

## Example
### Hello World
* Copy
* `go mod init ...`
* `go mod tidy`
* `go run main.go`
```golang
package main

import (
	"fmt"
	"os"

	"github.com/open-integration/core"
	"github.com/open-integration/core/pkg/state"
)

func main() {
	pipe := core.Pipeline{
		Metadata: core.PipelineMetadata{
			Name: "hello-world",
		},
		Spec: core.PipelineSpec{
			Reactions: []core.EventReaction{
				core.EventReaction{
					Condition: core.ConditionEngineStarted,
					Reaction: func(ev state.Event, state state.State) []core.Task {
						fmt.Println("Hello world")
						return []core.Task{}
					},
				},
			},
		},
	}
	e := core.NewEngine(&core.EngineOptions{
		Pipeline: pipe,
	})
	err := e.Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

```

### Real world examples
* [JIRA](https://github.com/olegsu/jira-sync)
* [Trello](https://github.com/olegsu/trello-sync)
