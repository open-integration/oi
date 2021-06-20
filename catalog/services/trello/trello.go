package trello

import (
	"context"
	_ "embed"

	"github.com/open-integration/oi/catalog/services/trello/endpoints/addcard"
	"github.com/open-integration/oi/catalog/services/trello/endpoints/archivecards"
	"github.com/open-integration/oi/catalog/services/trello/endpoints/getcards"
	"github.com/open-integration/oi/pkg/service"
	"github.com/open-integration/oi/pkg/utils"
)

//go:embed configs/endpoints/archivecard/arguments.json
var endpointArchivecardsArgument string

//go:embed configs/endpoints/archivecard/returns.json
var endpointArchivecardsReturns string

//go:embed configs/endpoints/addcard/arguments.json
var endpointAddcardArgument string

//go:embed configs/endpoints/addcard/returns.json
var endpointAddcardReturns string

//go:embed configs/endpoints/getcards/arguments.json
var endpointGetcardsArgument string

//go:embed configs/endpoints/getcards/returns.json
var endpointGetcardsReturns string

func Run(port string) {
	svc := service.New(port)
	err := svc.RegisterEndpoint("archivecards", archivecards.Handle, endpointArchivecardsArgument, endpointArchivecardsReturns)
	utils.DieOnError(err, "failed to register endpoint")

	err = svc.RegisterEndpoint("getcards", getcards.Handle, endpointGetcardsArgument, endpointGetcardsReturns)
	utils.DieOnError(err, "failed to register endpoint")

	err = svc.RegisterEndpoint("addcard", addcard.Handle, endpointAddcardArgument, endpointAddcardReturns)
	utils.DieOnError(err, "failed to register endpoint")

	err = svc.Run(context.Background())
	utils.DieOnError(err, "failed to run service")
}
