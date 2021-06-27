package message

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/open-integration/oi/catalog/services/slack/types"
	api "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/service"
)

func Handle(ctx context.Context, lgr logger.Logger, svc *service.Service, req *api.CallRequest) (*api.CallResponse, error) {
	args := &types.Arguments{}
	if err := service.UnmarshalRequestArgumentsInto(req, args); err != nil {
		return nil, err
	}
	if args.Message == nil {
		return service.BuildErrorResponse(fmt.Errorf("message is required"))
	}
	if args.WebhookURL == nil {
		return service.BuildErrorResponse(fmt.Errorf("webhook url is required"))
	}
	var buffer bytes.Buffer

	buffer.WriteString(`{ "text": "`)
	buffer.WriteString(*args.Message)
	buffer.WriteString(`"}`)

	lgr.Info("Sending message", "url", args.WebhookURL)

	res, err := http.Post(*args.WebhookURL, "application/x-www-form-urlencoded", strings.NewReader(buffer.String()))
	if err != nil {
		return service.BuildErrorResponse(err)
	}
	defer res.Body.Close()

	return service.BuildSuccessfullResponse(nil)
}
