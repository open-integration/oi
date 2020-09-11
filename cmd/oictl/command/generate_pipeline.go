package command

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/open-integration/oi/cmd/oictl/command/templates"
	"github.com/open-integration/oi/pkg/exec"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/template"
	"github.com/spf13/cobra"
)

var getwd = os.Getwd
var stat = os.Stat

type (
	rootGeneratePipelineCmdOptions struct {
		name                        string
		services                    []string
		engineStartCommand          string
		directory                   string
		skipGolangProjectInitiation bool
	}
)

var rootGeneratePipelineOptions rootGeneratePipelineCmdOptions

var rootGeneratePipelineCmd = &cobra.Command{
	Use:  "pipeline NAME [FLAGS] INTO_DIRECTORY",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootGeneratePipelineOptions.name = args[0]
		if len(args) == 2 {
			rootGeneratePipelineOptions.directory = args[1]
		}
		return execrootGeneratePipeline(&rootGeneratePipelineOptions)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		rootGenerateCmd.PreRun(cmd, args)
	},
}

func init() {
	rootGeneratePipelineCmd.PersistentFlags().StringVar(&rootGeneratePipelineOptions.engineStartCommand, "init-cmd", "", "Command to run when pipeline is started")
	rootGeneratePipelineCmd.PersistentFlags().BoolVar(&rootGeneratePipelineOptions.skipGolangProjectInitiation, "skip-project-init", false, "Do not create/update go.mod and go.sum files")
	rootGeneratePipelineCmd.PersistentFlags().StringArrayVar(&rootGeneratePipelineOptions.services, "service", []string{}, "services to use in format (name:version@as)")
	rootGenerateCmd.AddCommand(rootGeneratePipelineCmd)
}

func execrootGeneratePipeline(options *rootGeneratePipelineCmdOptions) error {
	log := logger.New(&logger.Options{
		LogToStdOut: rootOptions.verbose,
	})
	log.Debug("Generating pipeline", "name", options.name)
	services := []PipelineService{}
	for _, svc := range options.services {
		name, version, alias, err := splitServiceIntoParts(svc)
		if err != nil {
			return fmt.Errorf("Faild to read service format: %w", err)
		}
		log.Debug("Injecting service", "name", name, "alias", alias, "version", version)
		services = append(services, PipelineService{
			Name:    name,
			Version: version,
			As:      alias,
		})
	}

	log.Debug("Injecting task when pipeline is starting", "cmd", options.engineStartCommand)
	pipeline := &Pipeline{
		Metadata: Metadata{
			Name: options.name,
		},
		Services: services,
		EventReactions: []EventReaction{
			{
				Condition: "oi.ConditionEngineStarted()",
				Reaction: []Reaction{
					{
						Name:    "Command",
						Command: options.engineStartCommand,
					},
				},
			},
		},
	}
	res, err := template.Exec("pipeline", templates.PipelineTemplate, pipeline, nil)
	if err != nil {
		return fmt.Errorf("Failed to template pipeline: %w", err)
	}
	var mainFile io.Writer
	projectLocation := ""
	if options.directory == "" {
		mainFile = os.Stdout
	} else {
		location, err := resolveProjectFinalLocation(options.directory)
		if err != nil {
			return fmt.Errorf("Failed to resolved target directory: %w", err)
		}
		if err := ensureProjectLocation(location); err != nil {
			return fmt.Errorf("Failed to create project directory: %w", err)
		}
		projectLocation = location
		mainFile, err = os.Create(path.Join(location, "main.go"))
		if err != nil {
			return fmt.Errorf("Failed to create main.go: %w", err)
		}
	}
	if _, err := fmt.Fprintln(mainFile, res); err != nil {
		return fmt.Errorf("Failed to write pipeline to main.go: %w", err)
	}
	if !options.skipGolangProjectInitiation && projectLocation != "" {
		if err := exec.Exec(exec.Options{
			Command:    fmt.Sprintf("go mod init %s", pipeline.Metadata.Name),
			File:       os.Stdout,
			WorkingDir: projectLocation,
		}); err != nil {
			return fmt.Errorf("Failed to initiate Golang project with 'go mod init ... ': %w", err)
		}
		if err := exec.Exec(exec.Options{
			Command:    "go mod tidy",
			File:       os.Stdout,
			WorkingDir: projectLocation,
		}); err != nil {
			return fmt.Errorf("Failed to run 'go mod tidy': %w", err)
		}
	}
	if projectLocation != "" {
		if err := exec.Exec(exec.Options{
			Command:    "gofmt -l -w .",
			File:       os.Stdout,
			WorkingDir: projectLocation,
		}); err != nil {
			return fmt.Errorf("Failed to run 'gofmt -l -w .': %w", err)
		}
	}
	return nil
}

// splitServiceIntoParts returns service-name,version,alias
// from format NAME:VERSION@ALIAS
func splitServiceIntoParts(svc string) (string, string, string, error) {
	var name, version, alias string
	fullnameAlias := strings.Split(svc, "@")
	nameVersion := strings.Split(fullnameAlias[0], ":")
	name = nameVersion[0]
	if name == "" {
		return "", "", "", fmt.Errorf("Name must be part of service")
	}
	if len(nameVersion) == 2 {
		version = nameVersion[1]
	} else {
		version = "0.0.1"
	}

	_, err := semver.Parse(version)
	if err != nil {
		return "", "", "", fmt.Errorf("Version must be semantic version: %w", err)
	}

	if len(fullnameAlias) == 2 {
		alias = fullnameAlias[1]
	} else {
		alias = name
	}
	return name, version, alias, nil
}
