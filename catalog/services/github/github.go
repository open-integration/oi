package github

import (
	"context"

	"github.com/open-integration/oi/catalog/services/github/endpoints/getissuecomments"
	"github.com/open-integration/oi/catalog/services/github/endpoints/issuesearch"
	"github.com/open-integration/oi/pkg/service"
	"github.com/open-integration/oi/pkg/utils"
)

func Run(port string) {
	svc := service.New(port)
	err := svc.RegisterEndpoint("issuesearch", issuesearch.Handle)
	utils.DieOnError(err, "failed to register endpoint")
	err = svc.RegisterEndpoint("getissuecomments", getissuecomments.Handle)
	utils.DieOnError(err, "failed to register endpoint")
	err = svc.Run(context.Background())
	utils.DieOnError(err, "failed to run service")
}
