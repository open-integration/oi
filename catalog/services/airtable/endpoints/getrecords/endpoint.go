package getrecords

import (
	"context"

	at "github.com/mehanizm/airtable"
	"github.com/open-integration/oi/catalog/shared/airtable"
	api "github.com/open-integration/oi/pkg/api/v1"
	"github.com/open-integration/oi/pkg/logger"
	"github.com/open-integration/oi/pkg/service"
)

func Handle(ctx context.Context, lgr logger.Logger, svc *service.Service, req *api.CallRequest) (*api.CallResponse, error) {
	args := &GetRecordsArguments{}
	if err := service.UnmarshalRequestArgumentsInto(req, args); err != nil {
		return service.BuildErrorResponse(err)
	}

	table := airtable.GetTable(args.APIKey, args.DatabaseID, args.TableName)

	records, err := getRecords(table, lgr, args.Formula)
	if err != nil {
		return service.BuildErrorResponse(err)
	}

	return service.BuildSuccessfullResponse(GetRecordsReturnsClass{records})
}

func getRecords(table *at.Table, lgr logger.Logger, formula *string) ([]Record, error) {
	res := []Record{}
	request := table.GetRecords()
	if formula != nil {
		request.WithFilterFormula(*formula)
	}
	response, err := request.Do()
	if err != nil {
		return nil, err
	}
	for _, r := range response.Records {
		res = append(res, Record{
			Fields:      r.Fields,
			CreatedTime: &r.CreatedTime,
			Deleted:     &r.Deleted,
			ID:          &r.ID,
			Typecast:    &r.Typecast,
		})
	}

	lastOffset := response.Offset
	for {
		if lastOffset == "" {
			break
		}
		lgr.Info("use offset", "offest", lastOffset, "records", len(res))
		resp, err := table.GetRecords().WithOffset(lastOffset).Do()
		if err != nil {
			return nil, err
		}
		lastOffset = resp.Offset
		for _, r := range resp.Records {
			res = append(res, Record{
				Fields:      r.Fields,
				CreatedTime: &r.CreatedTime,
				Deleted:     &r.Deleted,
				ID:          &r.ID,
				Typecast:    &r.Typecast,
			})
		}
	}

	return res, nil
}
