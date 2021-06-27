package airtable

import (
	"context"

	"github.com/open-integration/oi/catalog/services/airtable/endpoints/addrecords"
	"github.com/open-integration/oi/catalog/services/airtable/endpoints/getrecords"
	"github.com/open-integration/oi/pkg/service"
	"github.com/open-integration/oi/pkg/utils"
)

func Run(port string) {
	svc := service.New(port)
	err := svc.RegisterEndpoint("getrecords", getrecords.Handle)
	utils.DieOnError(err, "failed to register endpoint")

	err = svc.RegisterEndpoint("addrecords", addrecords.Handle)
	utils.DieOnError(err, "failed to register endpoint")

	err = svc.Run(context.Background())
	utils.DieOnError(err, "failed to run service")
}
