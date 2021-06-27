package issuesearch

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/v35/github"
	"github.com/open-integration/oi/catalog/services/github/types"
	api "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/service"
	"golang.org/x/oauth2"
)

func Handle(ctx context.Context, lgr logger.Logger, svc *service.Service, req *api.CallRequest) (*api.CallResponse, error) {
	args := &types.IssueSearchArguments{}
	if err := service.UnmarshalRequestArgumentsInto(req, args); err != nil {
		return nil, err
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: args.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	res, _, err := client.Search.Issues(ctx, args.Query, nil)
	if err != nil {
		return service.BuildErrorResponse(err)
	}
	issues := &[]types.Issue{}
	b, err := json.Marshal(res.Issues)
	if err != nil {
		return service.BuildErrorResponse(err)
	}
	if err := json.Unmarshal(b, issues); err != nil {
		return service.BuildErrorResponse(err)
	}
	return service.BuildSuccessfullResponse(types.IssueSearchReturns{
		Issues: *issues,
	})
}
