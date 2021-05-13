package templates

// ServiceTemplateDockerfile template
var ServiceTemplateDockerfile = `
FROM alpine:3.9

RUN apk update && apk add --no-cache ca-certificates && apk upgrade

ARG NAME
ARG DIR

COPY ./$DIR/dist/$NAME .

RUN mv $NAME service

EXPOSE 80

ENV PORT=80

CMD "/service"
`

// ServiceTemplateMakefile template
var ServiceTemplateMakefile = `
.PHONY: build
build:
	@sh ./scripts/build.sh {{ .name }}
`

// ServiceTemplateVERSION template
var ServiceTemplateVERSION = `
{{ default "0.0.1" .version }}
`

// ServiceTemplateBuildScript template
var ServiceTemplateBuildScript = `
#!/bin/bash
set -e

build () {
    NAME=$1
    VERSION=$(cat ./VERSION)

    echo "Building $NAME-$VERSION service for darwin and amd64"
    GOOS=darwin GOARCH=amd64 packr build -o ./dist/$NAME-$VERSION-darwin-amd64 .

    echo "Building $NAME-$VERSION service for linux and amd64"
    GOOS=linux GOARCH=amd64 packr build -o ./dist/$NAME-$VERSION-linux-amd64 .

    echo "Building $NAME-$VERSION service for linux and 386"
    GOOS=linux GOARCH=386 packr build -o ./dist/$NAME-$VERSION-linux-386 .
}

build $1
`

// ServiceTemplateDefaultJSONSchema template
var ServiceTemplateDefaultJSONSchema = `
{
    "$id": "https://example.com/person.schema.json",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object"
}
`

// ServiceTemplateEndpoint template
var ServiceTemplateEndpoint = `
package {{ .endpoint.name }}

import (
	"context"
	"github.com/open-integration/oi/pkg/logger"
)

type (
	{{ .endpoint.name | strings.Title }}Options struct {
		Context   context.Context
		LoggerFD  string
		Arguments *{{ .endpoint.name | strings.Title }}Arguments
	}
)

func {{ .endpoint.name | strings.Title }}(opt {{ .endpoint.name | strings.Title }}Options) (*{{ .endpoint.name | strings.Title }}Returns, error) {
	log := logger.New(&logger.Options{
		FilePath: opt.LoggerFD,
	})
	log.Debug("dummy log output")
	return &{{ .endpoint.name | strings.Title }}Returns{}, nil
}
`

// ServiceTemplateMain template
var ServiceTemplateMain = `
package main

import (
	"context"
	"fmt"

	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"github.com/gobuffalo/packr"

	"github.com/open-integration/oi/pkg/logger"

	api "github.com/open-integration/oi/pkg/api/v1"
	{{ range .endpoints }}
	"{{ $.project }}/{{ $.name }}/pkg/endpoints/{{ .name }}"
	{{ end }}
)

type (
	Service struct {
		logger logger.Logger
		box    packr.Box
	}
)

func main() {

	service := &Service{
		logger: logger.New(&logger.Options{
			LogToStdOut: true,
		}),
		box:    packr.NewBox("./configs"),
	}
	runServer(context.Background(), service, os.Getenv("PORT"), service.logger)
}

func (s *Service) Init(context context.Context, req *api.InitRequest) (*api.InitResponse, error) {
	schemas := map[string]string{}
	for _, v := range s.box.List() {
		schema, _ := s.box.FindString(v)
		schemas[v] = schema
	}
	return &api.InitResponse{
		JsonSchemas: schemas,
	}, nil
}

func (s *Service) Call(context context.Context, req *api.CallRequest) (*api.CallResponse, error) {
	s.logger.Debug("Request", "endpoint", req.Endpoint)

	response := &api.CallResponse{}

	switch req.Endpoint {
	{{ range .endpoints }}
	case "{{ .name }}":
		args, err := {{ .name }}.Unmarshal{{ .name | strings.Title }}Arguments([]byte(req.Arguments))
		if resp := buildErrorResponse(err); resp != nil {
			return resp, nil
		}
		opt := {{ .name }}.{{ .name | strings.Title }}Options{
			Context:   context,
			LoggerFD:  req.Fd,
			Arguments: &args,
		}
		res, err := {{ .name }}.{{ .name | strings.Title }}(opt)
		if resp := buildErrorResponse(err); resp != nil {
			return resp, nil
		}
		payload, err := res.Marshal()
		if resp := buildErrorResponse(err); resp != nil {
			return resp, nil
		}
		response.Status = api.Status_OK
		response.Payload = string(payload)
		return response, nil
	{{ end }}
	}
	return buildErrorResponse(fmt.Errorf("Endpoint %s not found", req.Endpoint)), nil
}

func runServer(ctx context.Context, v1API api.ServiceServer, port string, log logger.Logger) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	api.RegisterServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Debug("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Debug("starting gRPC server", "port", port)
	err = server.Serve(listen)
	if err != nil {
		log.Debug("Error starting gRPC server", "error", err.Error())
		os.Exit(1)
	}
	return nil
}

func buildErrorResponse(err error) *api.CallResponse {
	if err != nil {
		return &api.CallResponse{
			Error:  err.Error(),
			Status: api.Status_Error,
		}
	}
	return nil
}
`
