package http

import (
	"context"
	_ "embed"

	"github.com/open-integration/oi/catalog/services/http/endpoints/call"
	"github.com/open-integration/oi/pkg/service"
)

//go:embed configs/endpoints/call/arguments.json
var endpointCallArgument string

//go:embed configs/endpoints/call/returns.json
var endpointCallReturns string

func Run(port string) {
	svc := service.New(port)
	svc.RegisterEndpoint("call", call.Handle, endpointCallArgument, endpointCallReturns)
	svc.Run(context.Background())
}
