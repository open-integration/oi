package getevents

import (
	"context"
	"encoding/json"

	"github.com/open-integration/oi/catalog/services/google-calendar/types"
	api "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/service"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func Handle(ctx context.Context, lgr logger.Logger, svc *service.Service, req *api.CallRequest) (*api.CallResponse, error) {
	args := &types.GoogleCalendarEventsArgumentsSchema{}
	if err := service.UnmarshalRequestArgumentsInto(req, args); err != nil {
		return nil, err
	}

	calendar, err := connect(args.ServiceAccount, lgr)
	if err != nil {
		return service.BuildErrorResponse(err)
	}
	call := calendar.Events.List(args.CalendarID)
	call.Context(ctx)
	call.ShowDeleted(*args.ShowDeleted)
	call.TimeMin(*args.TimeMin)
	call.TimeMax(*args.TimeMax)
	call.SingleEvents(true)

	events, err := call.Do()
	if err != nil {
		return service.BuildErrorResponse(err)
	}

	return service.BuildSuccessfullResponse(events.Items)
}

func connect(serviceAccount types.ServiceAccount, lgr logger.Logger) (*calendar.Service, error) {
	b, err := json.Marshal(serviceAccount)

	scopes := []string{
		calendar.CalendarScope,
		calendar.CalendarEventsReadonlyScope,
	}

	config, err := google.JWTConfigFromJSON(b, scopes...)
	if err != nil {
		return nil, err
	}
	client := config.Client(context.Background())
	return calendar.New(client)
}
