package command

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func dieOnError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

// resolveProjectFinalLocation returns the final location of the pipeline dir
func resolveProjectFinalLocation(target string) (string, error) {
	location := ""
	if strings.HasPrefix(target, ".") {
		dir, err := getwd()
		if err != nil {
			return "", fmt.Errorf("Failed to get current working directory: %w", err)
		}
		location = path.Join(dir, target)
	} else if strings.HasPrefix(target, "/") {
		location = target
	}
	return location, nil
}

// ensureProjectLocation creates project directory if not exist
func ensureProjectLocation(target string) error {
	if _, err := stat(target); os.IsNotExist(err) {
		if e := os.MkdirAll(target, os.ModePerm); e != nil {
			return fmt.Errorf("Failed to create project directory: %w", e)
		}
	}
	s, err := stat(target)
	if err != nil {
		return fmt.Errorf("Failed to get stat on project location %s: %w", target, err)
	}
	if !s.IsDir() {
		return fmt.Errorf("Project location %s is not directory", target)
	}
	return nil
}

// ensureTargetFile creates new file
func ensureTargetFile(defaultStd io.Writer, root string, location string, name string) (io.Writer, error) {
	if root == "" {
		return defaultStd, nil
	}
	if err := ensureProjectLocation(path.Join(root, location)); err != nil {
		return nil, err
	}
	return os.Create(path.Join(root, location, name))
}
