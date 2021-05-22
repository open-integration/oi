package github

import (
	"context"
	_ "embed"

	"github.com/open-integration/oi/catalog/services/github/endpoints/getissuecomments"
	"github.com/open-integration/oi/catalog/services/github/endpoints/issuesearch"
	"github.com/open-integration/oi/pkg/service"
	"github.com/open-integration/oi/pkg/utils"
)

//go:embed configs/endpoints/issuesearch/arguments.json
var endpointIssueSearchArgument string

//go:embed configs/endpoints/issuesearch/returns.json
var endpointIssueSearchReturns string

//go:embed configs/endpoints/getissuecomments/arguments.json
var endpointGetIssueCommentsArgument string

//go:embed configs/endpoints/getissuecomments/returns.json
var endpointGetIssueCommentsReturns string

func Run(port string) {
	svc := service.New(port)

	err := svc.RegisterEndpoint("issuesearch", issuesearch.Handle, endpointIssueSearchArgument, endpointIssueSearchReturns)
	utils.DieOnError(err, "failed to register endpoint")
	err = svc.RegisterEndpoint("getissuecomments", getissuecomments.Handle, endpointGetIssueCommentsArgument, endpointGetIssueCommentsReturns)
	utils.DieOnError(err, "failed to register endpoint")

	err = svc.Run(context.Background())
	utils.DieOnError(err, "failed to run service")
}
