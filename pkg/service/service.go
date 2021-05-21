package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"

	api "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
	"google.golang.org/grpc"
)

type (
	// Service implements the gRPC service to be used by the different catalog services
	Service struct {
		Logger    logger.Logger
		Port      string
		endpoints map[string]endpoint
	}

	endpoint struct {
		arguments string
		returns   string
		handler   endpointHandler
	}

	endpointHandler func(ctx context.Context, lgr logger.Logger, svc *Service, req *api.CallRequest) (*api.CallResponse, error)
)

func New(port string) Service {
	return Service{
		Logger: logger.New(&logger.Options{
			LogToStdOut: true,
		}),
		Port:      port,
		endpoints: map[string]endpoint{},
	}
}

func (s *Service) Init(ctx context.Context, req *api.InitRequest) (*api.InitResponse, error) {
	schemas := map[string]string{}
	for k, v := range s.endpoints {
		schemas[fmt.Sprintf("endpoints/%s/arguments.json", k)] = v.arguments
		schemas[fmt.Sprintf("endpoints/%s/returns.json", k)] = v.returns
	}
	return &api.InitResponse{
		JsonSchemas: schemas,
	}, nil
}

func (s *Service) Call(ctx context.Context, req *api.CallRequest) (*api.CallResponse, error) {
	s.Logger.Info("Request", "endpoint", req.Endpoint)
	ep, ok := s.endpoints[req.Endpoint]
	if !ok {
		return s.buildErrorResponse(fmt.Errorf("endpoint %s not found", req.Endpoint)), nil
	}
	lgr := logger.New(&logger.Options{
		FilePath: req.GetFd(),
	})
	resp, err := ep.handler(ctx, lgr, s, req)
	if err != nil {
		return s.buildErrorResponse(err), nil
	}
	if resp == nil {
		return s.buildErrorResponse(fmt.Errorf("failed to handle request")), nil
	}
	return resp, nil
}

func (s *Service) Run(ctx context.Context) error {
	listen, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	api.RegisterServiceServer(server, s)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			s.Logger.Debug("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	s.Logger.Info("starting gRPC server", "port", s.Port)
	err = server.Serve(listen)
	if err != nil {
		s.Logger.Info("Error starting gRPC server", "error", err.Error())
		os.Exit(1)
	}
	return nil
}

func (s *Service) RegisterEndpoint(name string, handler endpointHandler, argumentsSchema string, returnsSchems string) error {
	if _, found := s.endpoints[name]; !found {
		s.endpoints[name] = endpoint{
			arguments: argumentsSchema,
			returns:   returnsSchems,
			handler:   handler,
		}
		return nil
	}
	return fmt.Errorf("already exist")
}

func UnmarshalRequestArgumentsInto(req *api.CallRequest, into interface{}) error {
	args := req.GetArguments()
	return json.Unmarshal([]byte(args), into)
}

func BuildSuccessfullResponse(payload interface{}) (*api.CallResponse, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return BuildErrorResponse(err)
	}
	return &api.CallResponse{
		Status:  api.Status_OK,
		Payload: string(b),
	}, nil
}

func BuildErrorResponse(err error) (*api.CallResponse, error) {
	return &api.CallResponse{
		Status: api.Status_Error,
		Error:  err.Error(),
	}, nil
}

func (s *Service) buildErrorResponse(err error) *api.CallResponse {
	if err != nil {
		return &api.CallResponse{
			Error:  err.Error(),
			Status: api.Status_Error,
		}
	}
	return nil
}
