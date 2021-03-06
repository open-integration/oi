package command

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/open-integration/oi/cmd/oictl/command/templates"
	"github.com/open-integration/oi/pkg/exec"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/template"
	"github.com/spf13/cobra"
)

type (
	rootGenerateServiceCmdOptions struct {
		name                        string
		project                     string
		directory                   string
		types                       []string
		endpoints                   []string
		skipGolangProjectInitiation bool
		stdout                      io.Writer
	}

	flow struct {
		name         string
		render       []template.RenderTask
		preCommands  []string
		postCommands []string
	}
)

var rootGenerateServiceOptions rootGenerateServiceCmdOptions

var rootGenerateServiceCmd = &cobra.Command{
	Use:  "service NAME [FLAGS]",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootGenerateServiceOptions.name = args[0]
		rootGenerateServiceOptions.directory = fmt.Sprintf("./service-%s", args[0])
		rootGenerateServiceOptions.stdout = os.Stdout
		err := execrootGenerateService(rootGenerateServiceOptions)
		dieOnError(err)
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		rootGenerateCmd.PreRun(cmd, args)
	},
}

func init() {
	rootGenerateCmd.AddCommand(rootGenerateServiceCmd)
	rootGenerateServiceCmd.Flags().StringArrayVar(&rootGenerateServiceOptions.types, "type", []string{}, "Location to JSON-Schema type declarations, the name of the file will be used as type definition")
	rootGenerateServiceCmd.Flags().StringVar(&rootGenerateServiceOptions.project, "project", "", "Name of the project, default name will be set as the service name")
	rootGenerateServiceCmd.Flags().StringArrayVar(&rootGenerateServiceOptions.endpoints, "endpoints", []string{}, "Directory where arguments.json and returns.json can be found")
	rootGenerateServiceCmd.Flags().BoolVar(&rootGenerateServiceOptions.skipGolangProjectInitiation, "skip-project-init", false, "Do not create/update project related files")
}

func execrootGenerateService(options rootGenerateServiceCmdOptions) error {
	log := logger.New(&logger.Options{
		LogToStdOut: rootOptions.verbose,
	})

	if options.project == "" {
		options.project = options.name
	}

	projectLocation := ""
	{
		location, err := resolveProjectFinalLocation(options.directory)
		if err != nil {
			return fmt.Errorf("failed to resolved target directory: %w", err)
		}
		if err := ensureProjectLocation(location); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}
		projectLocation = location
	}
	log.Debug("Project location", "location", projectLocation)

	// Add hello-world endpoint
	if len(options.endpoints) == 0 {
		relative := path.Join("configs", "endpoints", "hello")
		argumentsFile, err := ensureTargetFile(options.stdout, projectLocation, relative, "arguments.json")
		if err != nil {
			return fmt.Errorf("failed to create hello-world arguments.json file: %w", err)
		}
		if err := template.Render(template.RenderTask{
			Content: templates.ServiceTemplateDefaultJSONSchema,
			Data:    nil,
			Name:    "arguments.json",
			Out:     argumentsFile,
		}); err != nil {
			return fmt.Errorf("failed to write hello-world arguments.json file: %w", err)
		}

		returnsFile, err := ensureTargetFile(options.stdout, projectLocation, relative, "returns.json")
		if err != nil {
			return fmt.Errorf("failed to create hello-world returns.json file: %w", err)
		}
		if err := template.Render(template.RenderTask{
			Content: templates.ServiceTemplateDefaultJSONSchema,
			Data:    nil,
			Name:    "returns.json",
			Out:     returnsFile,
		}); err != nil {
			return fmt.Errorf("failed to write hello-world returns.json file: %w", err)
		}
		options.endpoints = append(options.endpoints, path.Join(projectLocation, relative))
	}

	endpoints, err := buildEndpointsMap(options.endpoints)
	if err != nil {
		return err
	}

	rootScopedData := buildServiceData(&Service{
		Name:      options.name,
		Version:   "0.0.1",
		Endpoints: endpoints,
		Project:   options.name,
		Types:     options.types,
	}, log)
	log.Debug("Starting service generation")

	flows := []flow{}

	// generate main.go
	mainGoFileFlows, err := buildMainGoFlow(options.stdout, options.skipGolangProjectInitiation, projectLocation, options.name, options.name, rootScopedData)
	if err != nil {
		return err
	}
	flows = append(flows, mainGoFileFlows...)

	endpointFlows, err := buildFlows(endpoints, projectLocation, options.skipGolangProjectInitiation, options.stdout, rootScopedData)
	if err != nil {
		return err
	}
	flows = append(flows, endpointFlows...)
	flows = append(flows, flow{
		name: "Finalize project",
		postCommands: []string{
			"go mod tidy",
			"gofmt -l -w .",
		},
	})
	return runFlows(flows, projectLocation, options.stdout)
}

func buildServiceData(svc *Service, log logger.Logger) map[string]interface{} {
	data := map[string]interface{}{}
	{
		j, err := json.Marshal(svc)
		if err != nil {
			log.Error("failed to marshal service", "error", err.Error())
			return data
		}
		err = json.Unmarshal(j, &data)
		if err != nil {
			log.Error("failed to unmarshal service", "error", err.Error())
			return data
		}
	}
	return data
}

func buildEndpointData(endpoint Endpoint, root map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{}
	{
		ep := map[string]interface{}{}
		data["endpoint"] = ep
		dataJSONB, err := json.Marshal(endpoint)
		if err != nil {
			return data
		}
		err = json.Unmarshal(dataJSONB, &ep)
		if err != nil {
			return data
		}
		for k, v := range root {
			data[k] = v
		}
	}
	return data
}

func runPrePostCmds(cmds []string, workingDir string, out io.Writer) error {
	if workingDir == "" {
		return nil
	}
	for _, tmpl := range cmds {
		if err := exec.Exec(exec.Options{
			Command:    tmpl,
			WorkingDir: workingDir,
			File:       out,
		}); err != nil {
			return err
		}
	}
	return nil
}

func buildEndpointsMap(locations []string) (map[string]Endpoint, error) {
	endpoints := map[string]Endpoint{}
	for _, arg := range locations {
		name := path.Base(arg)
		_, ok := endpoints[name]
		if ok {
			return nil, fmt.Errorf("argument %s already exists as returns schema", name)
		}
		endpoints[name] = Endpoint{
			Name:      name,
			Arguments: path.Join(arg, "arguments.json"),
			Returns:   path.Join(arg, "returns.json"),
		}
	}
	return endpoints, nil
}

func buildFlows(endpoints map[string]Endpoint, projectLocation string, skipGolangProjectInitiation bool, stdout io.Writer, data map[string]interface{}) ([]flow, error) {
	flows := []flow{}
	for name, ep := range endpoints {
		epLocation := path.Join(projectLocation, "pkg", "endpoints", name)
		err := ensureProjectLocation(epLocation)
		if err != nil {
			return flows, fmt.Errorf("failed to create endpoint directory: %w", err)
		}
		post := []string{}
		if !skipGolangProjectInitiation {
			if ep.Arguments != "" {
				post = append(post, generateCodeGenerationCmd(epLocation, name, ep.Arguments, name, "arguments"))
			}
			if ep.Returns != "" {
				post = append(post, generateCodeGenerationCmd(epLocation, name, ep.Returns, name, "returns"))
			}
		}

		endpointFile, err := ensureTargetFile(stdout, projectLocation, path.Join("pkg", "endpoints", name), "endpoint.go")
		if err != nil {
			return flows, fmt.Errorf("failed to initiate main.go: %w", err)
		}

		flows = append(flows, flow{
			name: fmt.Sprintf("Generate endpoint %s", name),
			render: []template.RenderTask{
				{
					Content: templates.ServiceTemplateEndpoint,
					Data:    buildEndpointData(ep, data),
					Name:    "endpoint.go",
					Out:     endpointFile,
				},
			},
			postCommands: post,
		})
	}
	return flows, nil
}

func generateCodeGenerationCmd(epLocation string, epname string, sourceSchema string, golangPackage string, schemaType string) string {
	output := path.Join(epLocation, fmt.Sprintf("%s.go", schemaType))
	topName := fmt.Sprintf("%s%s", strings.Title(epname), strings.Title(schemaType))
	return fmt.Sprintf("quicktype -o %s -l go -s schema --src %s --package %s -t %s", output, sourceSchema, golangPackage, topName)
}

func buildMainGoFlow(stdout io.Writer, skipGolangProjectInitiation bool, projectLocation string, project string, name string, data map[string]interface{}) ([]flow, error) {
	flows := []flow{}
	pre := []string{}
	post := []string{}
	mainFile, err := ensureTargetFile(stdout, projectLocation, "", "main.go")
	if err != nil {
		return nil, fmt.Errorf("failed to initiate main.go: %w", err)
	}
	if !skipGolangProjectInitiation {
		pre = append(pre, fmt.Sprintf("go mod init %s/%s", project, name))
	}
	task := template.RenderTask{
		Content: templates.ServiceTemplateMain,
		Data:    data,
		Name:    "main.go",
		Out:     mainFile,
	}
	flows = append(flows, flow{
		name:         "Generate main.go file",
		render:       []template.RenderTask{task},
		preCommands:  pre,
		postCommands: post,
	})
	return flows, nil
}

func runFlows(flows []flow, projectLocation string, stdout io.Writer) error {
	var cmdoutFile io.Writer
	for _, flow := range flows {
		if projectLocation != "" {
			cmdoutFile = stdout
		}
		if err := runPrePostCmds(flow.preCommands, projectLocation, cmdoutFile); err != nil {
			return fmt.Errorf("failed to run pre commands: %w", err)
		}
		for _, task := range flow.render {
			if err := template.Render(task); err != nil {
				return fmt.Errorf("failed to render file %s: %w", task.Name, err)
			}
		}
		if err := runPrePostCmds(flow.postCommands, projectLocation, cmdoutFile); err != nil {
			return fmt.Errorf("failed to run post commands: %w", err)
		}
	}
	return nil
}
