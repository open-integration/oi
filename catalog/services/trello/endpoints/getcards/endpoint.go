package getcards

import (
	"context"
	"encoding/json"

	"github.com/adlio/trello"
	"github.com/open-integration/oi/catalog/services/trello/types"
	api "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/service"
)

func Handle(ctx context.Context, lgr logger.Logger, svc *service.Service, req *api.CallRequest) (*api.CallResponse, error) {
	args := &types.GetCardsArguments{}
	if err := service.UnmarshalRequestArgumentsInto(req, args); err != nil {
		return service.BuildErrorResponse(err)
	}
	lgr.Info("request", "app", args.Auth.App, "token", args.Auth.Token)
	client := trello.NewClient(args.Auth.App, args.Auth.Token)
	board, err := client.GetBoard(args.Board, trello.Defaults())
	if err != nil {
		return service.BuildErrorResponse(err)
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		return service.BuildErrorResponse(err)
	}
	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		return service.BuildErrorResponse(err)
	}

	for _, card := range cards {
		var list *trello.List
		for _, l := range lists {
			if card.IDList == l.ID {
				list = l
			}
		}
		card.List = list
	}
	j, err := json.Marshal(cards)
	if err != nil {
		return service.BuildErrorResponse(err)
	}
	res := []types.Card{}
	err = json.Unmarshal(j, &res)
	if err != nil {
		return service.BuildErrorResponse(err)
	}
	return service.BuildSuccessfullResponse(res)
}
