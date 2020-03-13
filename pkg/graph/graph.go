package graph

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
	"github.com/open-integration/core/pkg/state"
)

type (
	// Builder draw graphviz graph
	Builder interface {
		Build(state.State) []byte
	}

	builder struct{}
)

// New builds new graph drawer
func New() Builder {
	return &builder{}
}

func (b *builder) Build(s state.State) []byte {
	graph := gographviz.NewGraph()
	graph.SetDir(true)
	graph.Attrs.Add(string(gographviz.RankDir), "LR")
	for _, e := range s.Events() {
		if e.Metadata.Name == state.EventEngineStarted {
			graph.AddNode("G", "started", nil)
			for _, rt := range e.RelatedTasks {
				graph.AddNode("G", format(rt), node(s.Tasks()[rt].Status))
				graph.AddEdge("started", format(rt), true, nil)
			}
			continue
		}

		if len(e.RelatedTasks) > 0 {
			for _, rt := range e.RelatedTasks {
				graph.AddNode("G", format(rt), node(s.Tasks()[rt].Status))
				graph.AddEdge(format(e.Metadata.Task), format(rt), true, nil)
			}
		}
	}
	return []byte(graph.String())
}

func format(name string) string {
	return fmt.Sprintf("\"%s\"", name)
}

func statusToColor(status string) string {
	if status == state.TaskStatusSuccess {
		return "green"
	}
	if status == state.TaskStatusFailed {
		return "red"
	}
	return ""
}

func node(status string) map[string]string {
	return map[string]string{
		"color": statusToColor(status),
	}

}
