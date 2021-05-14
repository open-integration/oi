package task

import (
	"bytes"
	"fmt"
	"os"
	"path"
	coretemplate "text/template"

	"github.com/open-integration/oi/pkg/exec"

	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/template"
)

type (
	// RenderTask holds a structure to run a complete task
	// including renderring of templates
	// and running pre/post commands
	RenderTask struct {
		Name         string
		IDs          map[string]string
		Overwrite    bool
		PreCommands  []string
		PostCommands []string
		Data         interface{}
		Logger       logger.Logger
		Directory    string
		Func         coretemplate.FuncMap
	}
)

// Run RenderTask
func (r *RenderTask) Run() error {
	r.Logger.Debug("Starting task", "name", r.Name)

	r.Logger.Debug("Ensuring all ids exists in box")
	templates := map[string]string{}
	for name, tmpl := range r.IDs {
		filename, err := r.render(name, r.Data)
		if err != nil {
			return err
		}
		r.ensureDirectory(path.Dir(filename.String()))
		templates[filename.String()] = tmpl
	}

	if r.Directory != "" {
		r.ensureDirectory(r.Directory)
	}

	precmds := []string{}
	for _, pre := range r.PreCommands {
		res, err := r.render(pre, r.Data)
		if err != nil {
			return err
		}
		precmds = append(precmds, res.String())
	}
	if err := r.execCommands(precmds, r.Logger.New("state", "pre-commands")); err != nil {
		return err
	}

	err := r.writeTemplates(templates)
	if err != nil {
		return err
	}

	postcmds, err := r.buildPostCommands()
	if err != nil {
		return err
	}
	lgr := r.Logger.New("state", "post-commands")
	err = r.execCommands(postcmds, lgr)
	if err != nil {
		return err
	}

	return nil
}

func (r *RenderTask) runCommand(command string) error {
	return exec.Exec(exec.Options{
		Command: command,
	})
}

func (r *RenderTask) render(tmpl string, data interface{}) (*bytes.Buffer, error) {
	funcs := coretemplate.FuncMap{}
	for name, fn := range r.Func {
		funcs[name] = fn
	}
	return template.Exec("", tmpl, data, funcs)
}

func (r *RenderTask) write(location string, content *bytes.Buffer) error {
	r.Logger.Debug("Creating file", "location", location)
	f, err := os.Create(location)
	if err != nil {
		return err
	}
	r.Logger.Debug("File created")

	_, err = fmt.Fprintln(f, content)
	if err != nil {
		return err
	}
	r.Logger.Debug("File saved")
	return nil
}

func (r *RenderTask) ensureDirectory(dir string) {
	r.Logger.Debug("Ensuring target directory exist", "directory", dir)
	if !(isFileExist(dir)) && dir != "" {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			r.Logger.Error("Failed to create directory", "directory", dir)
		}
		r.Logger.Debug("Directory created", "directory", dir)
	}
}

func (r *RenderTask) execCommands(commands []string, logger logger.Logger) error {
	for _, cmd := range commands {
		logger.Debug("Running command", "cmd", cmd)
		if err := r.runCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

func (r *RenderTask) buildPostCommands() ([]string, error) {
	res := []string{}
	for _, post := range r.PostCommands {
		buf, err := r.render(post, r.Data)
		if err != nil {
			return nil, err
		}
		res = append(res, buf.String())
	}
	return res, nil
}

func (r *RenderTask) writeTemplates(templates map[string]string) error {
	for name, tmpl := range templates {
		shouldWrite := false
		filePath := path.Join(r.Directory, name)
		fileExist := isFileExist(filePath)
		r.Logger.Debug("Starting renderring templates", "name", name, "path", filePath)
		if !fileExist {
			shouldWrite = true
			r.Logger.Debug("File not exist", "file", filePath)
		} else if fileExist && r.Overwrite {
			r.Logger.Debug("File not exist, overwriting", "file", filePath)
			shouldWrite = true
		}

		if !shouldWrite {
			continue
		}

		data, err := r.render(tmpl, r.Data)
		if err != nil {
			return err
		}
		if err := r.write(filePath, data); err != nil {
			return err
		}
	}
	return nil
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false // file does not exist
	}
	return true
}
