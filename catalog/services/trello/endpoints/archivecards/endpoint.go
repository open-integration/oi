package archivecards

import (
	"context"

	"github.com/adlio/trello"
	api "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/service"
)

func Handle(ctx context.Context, lgr logger.Logger, svc *service.Service, req *api.CallRequest) (*api.CallResponse, error) {
	args := &ArchivecardsArguments{}
	if err := service.UnmarshalRequestArgumentsInto(req, args); err != nil {
		return service.BuildErrorResponse(err)
	}
	client := trello.NewClient(args.App, args.Token)

	for _, id := range args.CardIDs {
		if id == "" {
			continue
		}
		card, err := client.GetCard(id, trello.Defaults())
		if err != nil {
			lgr.Error("Failed to get card", "card", id, "error", err.Error())
			return service.BuildErrorResponse(err)
		}
		err = card.Update(trello.Arguments{
			"closed": "true",
		})
		if err != nil {
			lgr.Error("Failed to archive card", "card", id, "error", err.Error())
			continue
		}
		lgr.Debug("Card archived", "Card", card.ID)
	}
	return service.BuildSuccessfullResponse(&ArchivecardsReturns{})
}
