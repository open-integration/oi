package google_calendar

import (
	"context"

	"github.com/open-integration/oi/catalog/services/google-calendar/endpoints/getevents"
	"github.com/open-integration/oi/pkg/service"
	"github.com/open-integration/oi/pkg/utils"
)

func Run(port string) {
	svc := service.New(port)
	err := svc.RegisterEndpoint("getevents", getevents.Handle)
	utils.DieOnError(err, "failed to register endpoint")
	err = svc.Run(context.Background())
	utils.DieOnError(err, "failed to run service")
}
