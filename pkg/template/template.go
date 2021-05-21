package template

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/hairyhenderson/gomplate"
)

type (
	// RenderTask to render a template and print the output into target writer
	RenderTask struct {
		Name    string
		Content string
		Out     io.Writer
		Data    interface{}
	}
)

// Exec executes a template returns rendrred data
func Exec(name, tmpl string, data interface{}, funcs template.FuncMap) (*bytes.Buffer, error) {
	out := new(bytes.Buffer)
	f := gomplate.Funcs(nil)
	for name, fn := range funcs {
		f[name] = fn
	}
	t, err := template.New(name).Funcs(f).Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template string: %w", err)
	}
	if err := t.Execute(out, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	return out, nil
}

// Render renders a content and output into target io.Writer
func Render(t RenderTask) error {
	res, err := Exec(t.Name, t.Content, t.Data, nil)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}
	if _, err := fmt.Fprintln(t.Out, res); err != nil {
		return fmt.Errorf("failed to write data into target location: %w", err)
	}
	return nil
}
