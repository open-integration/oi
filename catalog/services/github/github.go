package github

import (
	"context"
	_ "embed"

	"github.com/open-integration/oi/catalog/services/github/endpoints/issuesearch"
	"github.com/open-integration/oi/pkg/service"
	"github.com/open-integration/oi/pkg/utils"
)

//go:embed configs/endpoints/issuesearch/arguments.json
var endpointIssueSeatchArgument string

//go:embed configs/endpoints/issuesearch/returns.json
var endpointIssueSeatchReturns string

func Run(port string) {
	svc := service.New(port)
	err := svc.RegisterEndpoint("issuesearch", issuesearch.Handle, endpointIssueSeatchArgument, endpointIssueSeatchReturns)
	utils.DieOnError(err, "failed to register endpoint")
	err = svc.Run(context.Background())
	utils.DieOnError(err, "failed to run service")
}
