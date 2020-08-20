package graph

import (
	"fmt"
	"time"

	"github.com/awalterschulze/gographviz"
	"github.com/open-integration/oi/core/state"
)

type (
	// Builder draw graphviz graph
	Builder interface {
		Build(state.State) ([]byte, error)
	}

	builder struct{}
)

// New builds new graph drawer
func New() Builder {
	return &builder{}
}

func (b *builder) Build(s state.State) ([]byte, error) {
	graph := gographviz.NewGraph()
	if err := graph.SetDir(true); err != nil {
		return nil, fmt.Errorf("Failed to set graph to be directed: %w", err)
	}
	if err := graph.Attrs.Add(string(gographviz.RankDir), "LR"); err != nil {
		return nil, fmt.Errorf("Failed to add graph LR attribute: %w", err)
	}
	for _, e := range s.Events() {
		if e.Metadata.Name == state.EventEngineStarted {
			if err := graph.AddNode("G", "started", nil); err != nil {
				return nil, fmt.Errorf("Failed to add first node to graph: %w", err)
			}
			for _, rt := range e.RelatedTasks {
				if err := graph.AddNode("G", format(rt), node(s.Tasks()[rt])); err != nil {
					return nil, fmt.Errorf("Failed to add node to graph: %w", err)
				}
				if err := graph.AddEdge("started", format(rt), true, nil); err != nil {
					return nil, fmt.Errorf("Failed to add edge to graph: %w", err)
				}
			}
			continue
		}

		if len(e.RelatedTasks) > 0 {
			for _, rt := range e.RelatedTasks {
				if err := graph.AddNode("G", format(rt), node(s.Tasks()[rt])); err != nil {
					return nil, fmt.Errorf("Failed to add node to graph: %w", err)
				}
				if err := graph.AddEdge(format(e.Metadata.Task), format(rt), true, nil); err != nil {
					return nil, fmt.Errorf("Failed to add edge to graph: %w", err)
				}
			}
		}
	}
	return []byte(graph.String()), nil
}

func format(name string) string {
	return fmt.Sprintf("\"%s\"", name)
}

func formatDiff(t time.Duration) string {
	return time.Time{}.Add(t).Format("15:04:05")
}

func statusToColor(status string) string {
	if status == state.TaskStatusSuccess {
		return "green"
	}
	if status == state.TaskStatusFailed {
		return "red"
	}
	return "\"\""
}

func node(taskstate state.TaskState) map[string]string {
	name := taskstate.Task.Name()
	times := taskstate.Times
	return map[string]string{
		"color": statusToColor(taskstate.Status),
		"label": fmt.Sprintf("\"{ <start> | <%s> name:%s \\n total:%s  | <end> }\"", name, name, formatDiff(times.Finished.Sub(times.Started))),
		"shape": "record",
	}
}
