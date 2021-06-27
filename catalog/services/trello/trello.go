package trello

import (
	"context"

	"github.com/open-integration/oi/catalog/services/trello/endpoints/addcard"
	"github.com/open-integration/oi/catalog/services/trello/endpoints/archivecards"
	"github.com/open-integration/oi/catalog/services/trello/endpoints/getcards"
	"github.com/open-integration/oi/pkg/service"
	"github.com/open-integration/oi/pkg/utils"
)

func Run(port string) {
	svc := service.New(port)
	err := svc.RegisterEndpoint("archivecards", archivecards.Handle)
	utils.DieOnError(err, "failed to register endpoint")

	err = svc.RegisterEndpoint("getcards", getcards.Handle)
	utils.DieOnError(err, "failed to register endpoint")

	err = svc.RegisterEndpoint("addcard", addcard.Handle)
	utils.DieOnError(err, "failed to register endpoint")

	err = svc.Run(context.Background())
	utils.DieOnError(err, "failed to run service")
}
