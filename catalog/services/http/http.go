package http

import (
	"context"
	_ "embed"

	"github.com/open-integration/oi/catalog/services/http/endpoints/call"
	"github.com/open-integration/oi/pkg/service"
	"github.com/open-integration/oi/pkg/utils"
)

//go:embed configs/endpoints/call/arguments.json
var endpointCallArgument string

//go:embed configs/endpoints/call/returns.json
var endpointCallReturns string

func Run(port string) {
	svc := service.New(port)
	err := svc.RegisterEndpoint("call", call.Handle, endpointCallArgument, endpointCallReturns)
	utils.DieOnError(err, "failed to register endpoint")
	err = svc.Run(context.Background())
	utils.DieOnError(err, "failed to run service")
}
