package exec

import (
	"fmt"
	"io"
	"os/exec"
)

type (
	// Options to execute command
	Options struct {
		Command     string
		WorkingDir  string
		Environment []string
		File        io.Writer
	}
)

// Exec executes as subprocess
func Exec(opt Options) error {
	cmd := exec.Command("sh", "-c", opt.Command) // #nosec G204
	cmd.Dir = opt.WorkingDir
	cmd.Env = opt.Environment
	cmd.Stdout = opt.File
	cmd.Stderr = opt.File
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to executed command: %s: %w", opt.Command, err)
	}
	return nil
}
