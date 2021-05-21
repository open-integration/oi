package slack

import (
	"context"
	_ "embed"

	"github.com/open-integration/oi/catalog/services/slack/endpoints/message"
	"github.com/open-integration/oi/pkg/service"
	"github.com/open-integration/oi/pkg/utils"
)

//go:embed configs/endpoints/message/arguments.json
var endpointMessageArgument string

//go:embed configs/endpoints/message/returns.json
var endpointMessageReturns string

func Run(port string) {
	svc := service.New(port)
	err := svc.RegisterEndpoint("message", message.Handle, endpointMessageArgument, endpointMessageReturns)
	utils.DieOnError(err, "failed to register endpoint")
	err = svc.Run(context.Background())
	utils.DieOnError(err, "failed to run service")
}
