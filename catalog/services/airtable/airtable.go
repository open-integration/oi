package airtable

import (
	"context"
	_ "embed"

	"github.com/open-integration/oi/catalog/services/airtable/endpoints/addrecords"
	"github.com/open-integration/oi/catalog/services/airtable/endpoints/getrecords"
	"github.com/open-integration/oi/pkg/service"
	"github.com/open-integration/oi/pkg/utils"
)

//go:embed configs/endpoints/getrecords/arguments.json
var endpointGetrecordsArgument string

//go:embed configs/endpoints/getrecords/returns.json
var endpointGetrecordsReturns string

//go:embed configs/endpoints/addrecords/arguments.json
var endpointAddrecordsArgument string

//go:embed configs/endpoints/addrecords/returns.json
var endpointAddrecordsReturns string

func Run(port string) {
	svc := service.New(port)
	err := svc.RegisterEndpoint("getrecords", getrecords.Handle, endpointGetrecordsArgument, endpointGetrecordsReturns)
	utils.DieOnError(err, "failed to register endpoint")

	err = svc.RegisterEndpoint("addrecords", addrecords.Handle, endpointAddrecordsArgument, endpointAddrecordsReturns)
	utils.DieOnError(err, "failed to register endpoint")

	err = svc.Run(context.Background())
	utils.DieOnError(err, "failed to run service")
}
