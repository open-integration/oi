package addcard

import (
	"context"

	"github.com/adlio/trello"
	"github.com/open-integration/oi/catalog/services/trello/types"
	api "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/service"
)

func Handle(ctx context.Context, lgr logger.Logger, svc *service.Service, req *api.CallRequest) (*api.CallResponse, error) {
	args := &types.AddCardArguments{}
	if err := service.UnmarshalRequestArgumentsInto(req, args); err != nil {
		return service.BuildErrorResponse(err)
	}
	client := trello.NewClient(args.Auth.App, args.Auth.Token)
	card := &trello.Card{
		Name:     args.Name,
		IDList:   args.List,
		IDLabels: args.Labels,
	}
	if args.Description != nil {
		card.Desc = *args.Description
	}

	err := client.CreateCard(card, trello.Defaults())
	if err != nil {
		return service.BuildErrorResponse(err)
	}
	return service.BuildSuccessfullResponse(&types.AddCardReturns{})
}
